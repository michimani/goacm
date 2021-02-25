package gacm

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// MockACMDescribeCertificateAPI is a type that represents a function that mock ACM's DescribeCertificate.
type MockACMDescribeCertificateAPI func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)

// MockACMListCertificatesAPI is a type that represents a function that mock ACM's ListCertificatesOutput.
type MockACMListCertificatesAPI func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)

// DescribeCertificate returns a function that mock original of ACM DescribeCertificate.
func (m MockACMDescribeCertificateAPI) DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
	return m(ctx, params, optFns...)
}

// ListCertificates returns a function that mock original of ACM ListCertificates.
func (m MockACMListCertificatesAPI) ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return m(ctx, params, optFns...)
}

// GenerateMockACMDescribeCertificateAPI returns MockACMDescribeCertificateAPI
func GenerateMockACMDescribeCertificateAPI(t *testing.T, arn, domainName, status, cType, failurReason string) MockACMDescribeCertificateAPI {
	return MockACMDescribeCertificateAPI(func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
		t.Helper()
		if params.CertificateArn == nil {
			t.Fatal("expect Certificate ARN to not be nil")
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
