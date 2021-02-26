package acmgo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// ACMgo is a structure that wraps an ACM client.
type ACMgo struct {
	client *acm.Client
	region string
}

// NewACMgo returns a new ACMgo object.
func NewACMgo(region string) (*ACMgo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &ACMgo{
		client: acm.NewFromConfig(cfg),
		region: region,
	}, nil
}

// ListCertificateSummaries returns a list of certificate summary.
func ListCertificateSummaries(api ACMListCertificatesAPI) ([]types.CertificateSummary, error) {
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
		ARN:           arn,
		DomainName:    aws.ToString(out.Certificate.DomainName),
		Status:        out.Certificate.Status,
		Type:          out.Certificate.Type,
		FailureReason: out.Certificate.FailureReason,
	}, nil
}

// ListCertificates is returns list of certificate.
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
