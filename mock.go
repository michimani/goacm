package acmgo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// MockParams is a structure with the elements needed to generate a mock.
type MockParams struct {
	Arn             string
	DomainName      string
	ArnBase         string
	DomainNameBase  string
	Status          string
	CertificateType string
	FailureReason   string
	Count           int
}

// MockACMAPI is a struct that represents an ACM client.
type MockACMAPI struct {
	DescribeCertificateAPI MockACMDescribeCertificateAPI
	ListCertificatesAPI    MockACMListCertificatesAPI
}

// MockACMDescribeCertificateAPI is a type that represents a function that mock ACM's DescribeCertificate.
type MockACMDescribeCertificateAPI func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)

// MockACMListCertificatesAPI is a type that represents a function that mock ACM's ListCertificatesOutput.
type MockACMListCertificatesAPI func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)

// DescribeCertificate returns a function that mock original of ACM DescribeCertificate.
func (m MockACMAPI) DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
	return m.DescribeCertificateAPI(ctx, params, optFns...)
}

// ListCertificates returns a function that mock original of ACM ListCertificates.
func (m MockACMAPI) ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return m.ListCertificatesAPI(ctx, params, optFns...)
}

// GenerateMockACMAPI return MockACMAPI.
func GenerateMockACMAPI(p MockParams) MockACMAPI {
	return MockACMAPI{
		DescribeCertificateAPI: GenerateMockACMDescribeCertificateAPI(p.Arn, p.DomainName, p.Status, p.CertificateType, p.FailureReason),
		ListCertificatesAPI:    GenerateMockACMListCertificatesAPI(p.ArnBase, p.DomainNameBase, p.Count),
	}
}

// GenerateMockACMDescribeCertificateAPI returns MockACMDescribeCertificateAPI.
func GenerateMockACMDescribeCertificateAPI(arn, domainName, status, cType, failurReason string) MockACMDescribeCertificateAPI {
	return MockACMDescribeCertificateAPI(func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
		if params.CertificateArn == nil {
			return nil, errors.New("expect Certificate ARN to not be nil")
		}

		if aws.ToString(params.CertificateArn) != arn {
			return nil, errors.New("not found")
		}

		return &acm.DescribeCertificateOutput{
			Certificate: &types.CertificateDetail{
				CertificateArn: aws.String(arn),
				DomainName:     aws.String(domainName),
				Status:         types.CertificateStatus(status),
				Type:           types.CertificateType(cType),
				FailureReason:  types.FailureReason(failurReason),
			},
		}, nil
	})
}

// GenerateMockACMListCertificatesAPI returns MockACMDescribeCertificateAPI.
func GenerateMockACMListCertificatesAPI(arnBase, domainBase string, count int) MockACMListCertificatesAPI {
	return MockACMListCertificatesAPI(func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
		var csList []types.CertificateSummary

		for i := 0; i < count; i++ {
			csList = append(csList, types.CertificateSummary{
				CertificateArn: aws.String(arnBase + strconv.Itoa(i+1)),
				DomainName:     aws.String(fmt.Sprintf("test%s.%s", strconv.Itoa(i+1), domainBase)),
			})
		}

		return &acm.ListCertificatesOutput{
			CertificateSummaryList: csList,
		}, nil
	})
}
