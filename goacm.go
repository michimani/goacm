package goacm

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
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
func NewGoACM(region string) (*GoACM, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
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
func ListCertificateSummaries(api ACMListCertificatesAPI) ([]acmTypes.CertificateSummary, error) {
	in := acm.ListCertificatesInput{}
	out, err := api.ListCertificates(context.TODO(), &in)
	if err != nil {
		return nil, err
	}

	return out.CertificateSummaryList, nil
}

// GetCertificate returns the details of the certificate.
func GetCertificate(api ACMDescribeCertificateAPI, arn string) (Certificate, error) {
	in := acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	}
	out, err := api.DescribeCertificate(context.TODO(), &in)
	if err != nil {
		return Certificate{}, err
	}

	return Certificate{
		Arn:           arn,
		DomainName:    aws.ToString(out.Certificate.DomainName),
		Status:        out.Certificate.Status,
		Type:          out.Certificate.Type,
		FailureReason: out.Certificate.FailureReason,
	}, nil
}

// ListCertificates returns list of certificate.
func ListCertificates(api ACMAPI) ([]Certificate, error) {
	summary, err := ListCertificateSummaries(api)
	if err != nil {
		return nil, err
	}

	var cList []Certificate
	for _, s := range summary {
		c, err := GetCertificate(api, aws.ToString(s.CertificateArn))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		cList = append(cList, c)
	}

	return cList, nil
}

// DeleteCertificate returns an error if deleting the certificate fails.
func DeleteCertificate(api ACMDeleteCertificateAPI, arn string) error {
	in := acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	}

	if _, err := api.DeleteCertificate(context.TODO(), &in); err != nil {
		return err
	}

	return nil
}

// IssueCertificate issues an SSL certificate for the specified domain.
func IssueCertificate(aAPI ACMAPI, rAPI Route53API, method ValidationMethod, targetDomain, hostedDomain string) error {
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
	r, err := aAPI.RequestCertificate(context.TODO(), &reqIn)
	if err != nil {
		return err
	}

	if method == ValidationMethodEmail {
		return nil
	}

	dcIn := acm.DescribeCertificateInput{
		CertificateArn: r.CertificateArn,
	}
	c, err := aAPI.DescribeCertificate(context.TODO(), &dcIn)
	if err != nil {
		return err
	}
	if c.Certificate.DomainValidationOptions == nil {
		return errors.New("DomainValidationOptions dose not exists")
	}

	vRecordName := c.Certificate.DomainValidationOptions[0].ResourceRecord.Name
	vRecordValue := c.Certificate.DomainValidationOptions[0].ResourceRecord.Value

	lhzIn := route53.ListHostedZonesInput{}
	h, err := rAPI.ListHostedZones(context.TODO(), &lhzIn)
	if err != nil {
		return err
	}

	hzID := ""
	for _, hz := range h.HostedZones {
		if *hz.Name == hostedDomain {
			hzID = *hz.Id
		}
	}
	if hzID == "" {
		return errors.New("Cannot get hosted zone ID")
	}

	crsIn := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hzID),
		ChangeBatch: &route53Types.ChangeBatch{
			Changes: []route53Types.Change{
				{
					Action: route53Types.ChangeActionCreate,
					ResourceRecordSet: &route53Types.ResourceRecordSet{
						Name: vRecordName,
						Type: route53Types.RRType("CNAME"),
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

	_, err = rAPI.ChangeResourceRecordSets(context.TODO(), &crsIn)
	if err != nil {
		return err
	}

	return nil
}

// RollbackIssueCertificate rollbacks to issue an SSL certificate.
func RollbackIssueCertificate(aAPI ACMAPI, rAPI Route53API) error {
	return nil
}
