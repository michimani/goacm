package goacm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// GoACM is a structure that wraps an ACM client.
type GoACM struct {
	ACMClient     *acm.Client
	Route53Client *route53.Client
	Region        string
}

// NewGoACM returns a new GoACM object.
func NewGoACM(ctx context.Context, region string) (*GoACM, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &GoACM{
		ACMClient:     acm.NewFromConfig(cfg),
		Route53Client: route53.NewFromConfig(cfg),
		Region:        region,
	}, nil
}

// ListCertificateSummaries returns a list of certificate summary.
func ListCertificateSummaries(ctx context.Context, api ACMListCertificatesAPI) ([]acmTypes.CertificateSummary, error) {
	in := acm.ListCertificatesInput{}
	out, err := api.ListCertificates(ctx, &in)
	if err != nil {
		return nil, err
	}

	return out.CertificateSummaryList, nil
}

// GetCertificate returns the details of the certificate.
func GetCertificate(ctx context.Context, api ACMDescribeCertificateAPI, arn string) (Certificate, error) {
	in := acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	}
	out, err := api.DescribeCertificate(ctx, &in)
	if err != nil {
		return Certificate{}, err
	}

	vMethod := ""
	recordSet := RecordSet{}
	if out.Certificate.DomainValidationOptions != nil {
		vMethod = string(out.Certificate.DomainValidationOptions[0].ValidationMethod)
		if vMethod == string(types.ValidationMethodDns) {
			recordSet.HostedDomainName = aws.ToString(out.Certificate.DomainValidationOptions[0].ValidationDomain)
			recordSet.Name = aws.ToString(out.Certificate.DomainValidationOptions[0].ResourceRecord.Name)
			recordSet.Value = aws.ToString(out.Certificate.DomainValidationOptions[0].ResourceRecord.Value)
			recordSet.Type = string(out.Certificate.DomainValidationOptions[0].ResourceRecord.Type)
		}
	}

	return Certificate{
		Arn:                 arn,
		DomainName:          aws.ToString(out.Certificate.DomainName),
		Status:              string(out.Certificate.Status),
		Type:                string(out.Certificate.Type),
		FailureReason:       string(out.Certificate.FailureReason),
		ValidationMethod:    vMethod,
		ValidationRecordSet: recordSet,
	}, nil
}

// ListCertificates returns list of certificate.
func ListCertificates(ctx context.Context, api ACMAPI) ([]Certificate, error) {
	summary, err := ListCertificateSummaries(ctx, api)
	if err != nil {
		return nil, err
	}

	var cList []Certificate
	for _, s := range summary {
		c, err := GetCertificate(ctx, api, aws.ToString(s.CertificateArn))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		cList = append(cList, c)
	}

	return cList, nil
}

// DeleteCertificate returns an error if deleting the certificate fails.
func DeleteCertificate(ctx context.Context, aAPI ACMAPI, rAPI Route53API, arn string) error {
	c, err := GetCertificate(ctx, aAPI, arn)
	if err != nil {
		return err
	}

	// Delete Route 53 Record that validate domain.
	if c.ValidationMethod == string(types.ValidationMethodDns) {
		if err := DeleteRoute53RecordSet(ctx, aAPI, rAPI, c.ValidationRecordSet); err != nil {
			return err
		}
	}

	in := acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	}

	if _, err := aAPI.DeleteCertificate(ctx, &in); err != nil {
		return err
	}

	return nil
}

// IssueCertificate issues an SSL certificate for the specified domain.
func IssueCertificate(ctx context.Context, aAPI ACMAPI, rAPI Route53API, method, targetDomain, hostedDomain string) (IssueCertificateResult, error) {
	var result IssueCertificateResult = IssueCertificateResult{
		DomainName:       targetDomain,
		HostedDomainName: hostedDomain,
		ValidationMethod: string(method),
	}

	// request certificate
	reqIn := acm.RequestCertificateInput{
		DomainName:       aws.String(targetDomain),
		ValidationMethod: acmTypes.ValidationMethod(method),
		DomainValidationOptions: []acmTypes.DomainValidationOption{
			{
				DomainName:       aws.String(targetDomain),
				ValidationDomain: aws.String(hostedDomain),
			},
		},
	}
	r, err := aAPI.RequestCertificate(ctx, &reqIn)
	if err != nil {
		return IssueCertificateResult{}, err
	}

	result.CertificateArn = aws.ToString(r.CertificateArn)

	if method == string(types.ValidationMethodEmail) {
		return result, nil
	}

	time.Sleep(time.Second * 5)

	dcIn := acm.DescribeCertificateInput{
		CertificateArn: r.CertificateArn,
	}
	c, err := aAPI.DescribeCertificate(ctx, &dcIn)
	if err != nil {
		return IssueCertificateResult{}, err
	}
	if c.Certificate.DomainValidationOptions == nil {
		errMsg := "DomainValidationOptions dose not exists"
		if err := RollbackIssueCertificate(ctx, aAPI, rAPI, *c.Certificate.CertificateArn); err != nil {
			errMsg += fmt.Sprintf("; Failed to rollback to issue certificate: %v", err)
		} else {
			errMsg += "; rollbacked to issue certificate"
		}
		return IssueCertificateResult{}, errors.New(errMsg)
	}

	vRecordName := c.Certificate.DomainValidationOptions[0].ResourceRecord.Name
	vRecordValue := c.Certificate.DomainValidationOptions[0].ResourceRecord.Value

	result.ValidationRecordName = *vRecordName
	result.ValidationRecordValue = *vRecordValue

	// allowed only public hosted zones
	hzID, err := getPublicHostedZoneIDByDomainName(ctx, rAPI, hostedDomain)
	if err != nil {
		errMsg := err.Error()
		if err := RollbackIssueCertificate(ctx, aAPI, rAPI, *c.Certificate.CertificateArn); err != nil {
			errMsg += fmt.Sprintf("; Failed to rollback to issue certificate: %v", err)
		} else {
			errMsg += "; rollbacked to issue certificate"
		}
		return IssueCertificateResult{}, errors.New(errMsg)
	}

	if hzID == "" {
		errMsg := "Cannot get public hosted zone ID"
		if err := RollbackIssueCertificate(ctx, aAPI, rAPI, *c.Certificate.CertificateArn); err != nil {
			errMsg += fmt.Sprintf("; Failed to rollback to issue certificate: %v", err)
		} else {
			errMsg += "; rollbacked to issue certificate"
		}
		return IssueCertificateResult{}, errors.New(errMsg)
	}

	result.HosteZoneID = hzID

	crsIn := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hzID),
		ChangeBatch: &route53Types.ChangeBatch{
			Changes: []route53Types.Change{
				{
					Action: route53Types.ChangeActionCreate,
					ResourceRecordSet: &route53Types.ResourceRecordSet{
						Name: vRecordName,
						Type: route53Types.RRTypeCname,
						TTL:  aws.Int64(300),
						ResourceRecords: []route53Types.ResourceRecord{
							{
								Value: vRecordValue,
							},
						},
					},
				},
			},
		},
	}

	_, err = rAPI.ChangeResourceRecordSets(ctx, &crsIn)
	if err != nil {
		errMsg := err.Error()
		if err := RollbackIssueCertificate(ctx, aAPI, rAPI, *c.Certificate.CertificateArn); err != nil {
			errMsg += fmt.Sprintf("; Failed to rollback to issue certificate: %v", err)
		} else {
			errMsg += "; rollbacked to issue certificate"
		}
		return IssueCertificateResult{}, errors.New(errMsg)
	}

	return result, nil
}

// RollbackIssueCertificate rollbacks to issue an SSL certificate.
func RollbackIssueCertificate(ctx context.Context, aAPI ACMAPI, rAPI Route53API, arn string) error {
	return DeleteCertificate(ctx, aAPI, rAPI, arn)
}

// DeleteRoute53RecordSet deletes a Route 53 record set.
func DeleteRoute53RecordSet(ctx context.Context, aAPI ACMAPI, rAPI Route53API, rs RecordSet) error {
	hzID, err := getPublicHostedZoneIDByDomainName(ctx, rAPI, rs.HostedDomainName)
	if err != nil {
		return err
	}
	if hzID == "" {
		return errors.New("Cannot get hosted zone ID")
	}

	lrrsIn := route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hzID),
		StartRecordName: aws.String(rs.Name),
		MaxItems:        aws.Int32(1),
	}
	r, err := rAPI.ListResourceRecordSets(ctx, &lrrsIn)
	if err != nil {
		return err
	}

	if len(r.ResourceRecordSets) != 1 {
		return fmt.Errorf("Target RecordeSet does not exists: %s", rs.Name)
	}

	rrs := r.ResourceRecordSets[0]
	if aws.ToString(rrs.Name) != rs.Name {
		return fmt.Errorf("Target RecordeSet does not exists: %s", rs.Name)
	}

	crsIn := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hzID),
		ChangeBatch: &route53Types.ChangeBatch{
			Changes: []route53Types.Change{
				{
					Action: route53Types.ChangeActionDelete,
					ResourceRecordSet: &route53Types.ResourceRecordSet{
						Name: aws.String(rs.Name),
						Type: route53Types.RRType(rs.Type),
						TTL:  rrs.TTL,
						ResourceRecords: []route53Types.ResourceRecord{
							{
								Value: aws.String(rs.Value),
							},
						},
					},
				},
			},
		},
	}

	_, err = rAPI.ChangeResourceRecordSets(ctx, &crsIn)

	if err != nil {
		return err
	}

	return nil
}

// Get public hosted zone ID by domain name.
// domainName is a string without a "." at the end.
func getPublicHostedZoneIDByDomainName(ctx context.Context, rAPI Route53API, domainName string) (string, error) {
	dn := domainName + "."

	lhzIn := route53.ListHostedZonesInput{}
	out, err := rAPI.ListHostedZones(ctx, &lhzIn)
	if err != nil {
		return "", err
	}

	for _, hz := range out.HostedZones {
		if hz.Config != nil && !hz.Config.PrivateZone && *hz.Name == dn {
			return aws.ToString(hz.Id), nil
		}
	}

	return "", nil
}
