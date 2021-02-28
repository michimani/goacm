package goacm

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// MockParams is a structure with the elements needed to generate a mock.
type MockParams struct {
	Arn             string
	DomainName      string
	Status          string
	CertificateType string
	FailureReason   string
	Count           int
}

// MockACMAPI is a struct that represents an ACM client.
type MockACMAPI struct {
	ListCertificatesAPI    MockACMListCertificatesAPI
	DescribeCertificateAPI MockACMDescribeCertificateAPI
	DeleteCertificateAPI   MockACMDeleteCertificateAPI
	RequestCertificateAPI  MockACMRequestCertificateAPI
}

// MockACMDescribeCertificateAPI is a type that represents a function that mock ACM's DescribeCertificate.
type MockACMDescribeCertificateAPI func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)

// MockACMListCertificatesAPI is a type that represents a function that mock ACM's ListCertificates.
type MockACMListCertificatesAPI func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)

// MockACMDeleteCertificateAPI is a type that represents a function that mock ACM's DeleteCertificate.
type MockACMDeleteCertificateAPI func(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error)

// MockACMRequestCertificateAPI is a type that represents a function that mock ACM's RequestCertificate.
type MockACMRequestCertificateAPI func(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error)

// DescribeCertificate returns a function that mock original of ACM DescribeCertificate.
func (m MockACMAPI) DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
	return m.DescribeCertificateAPI(ctx, params, optFns...)
}

// ListCertificates returns a function that mock original of ACM ListCertificates.
func (m MockACMAPI) ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return m.ListCertificatesAPI(ctx, params, optFns...)
}

// DeleteCertificate returns a function that mock original of ACM DeleteCertificate.
func (m MockACMAPI) DeleteCertificate(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error) {
	return m.DeleteCertificateAPI(ctx, params, optFns...)
}

// RequestCertificate returns a function that mock original of ACM RequestCertificate.
func (m MockACMAPI) RequestCertificate(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error) {
	return m.RequestCertificate(ctx, params, optFns...)
}

// GenerateMockACMAPI return MockACMAPI.
func GenerateMockACMAPI(mockParams []MockParams) MockACMAPI {
	return MockACMAPI{
		DescribeCertificateAPI: GenerateMockACMDescribeCertificateAPI(mockParams),
		ListCertificatesAPI:    GenerateMockACMListCertificatesAPI(mockParams),
		DeleteCertificateAPI:   GenerateMockACMDeleteCertificateAPI(mockParams),
		RequestCertificateAPI:  GenerateMockACMRequestCertificateAPI(mockParams),
	}
}

// GenerateMockACMDescribeCertificateAPI returns MockACMDescribeCertificateAPI.
func GenerateMockACMDescribeCertificateAPI(mockParams []MockParams) MockACMDescribeCertificateAPI {
	return MockACMDescribeCertificateAPI(func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
		if params.CertificateArn == nil {
			return nil, errors.New("expect Certificate ARN to not be nil")
		}

		var availableCertificates map[string]*acm.DescribeCertificateOutput = map[string]*acm.DescribeCertificateOutput{}
		for _, mp := range mockParams {
			if mp.Arn == "" {
				continue
			}
			availableCertificates[mp.Arn] = &acm.DescribeCertificateOutput{
				Certificate: &types.CertificateDetail{
					CertificateArn: aws.String(mp.Arn),
					DomainName:     aws.String(mp.DomainName),
					Status:         types.CertificateStatus(mp.Status),
					Type:           types.CertificateType(mp.CertificateType),
					FailureReason:  types.FailureReason(mp.FailureReason),
				},
			}
		}

		dco := availableCertificates[*params.CertificateArn]
		if dco == nil {
			return nil, errors.New("not found")
		}

		return dco, nil
	})
}

// GenerateMockACMListCertificatesAPI returns MockACMDescribeCertificateAPI.
func GenerateMockACMListCertificatesAPI(mockParams []MockParams) MockACMListCertificatesAPI {
	return MockACMListCertificatesAPI(func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
		var csList []types.CertificateSummary

		for _, mp := range mockParams {
			csList = append(csList, types.CertificateSummary{
				CertificateArn: aws.String(mp.Arn),
				DomainName:     aws.String(mp.DomainName),
			})
		}

		return &acm.ListCertificatesOutput{
			CertificateSummaryList: csList,
		}, nil
	})
}

// GenerateMockACMDeleteCertificateAPI returns MockACMDeleteCertificateAPI
func GenerateMockACMDeleteCertificateAPI(mockParams []MockParams) MockACMDeleteCertificateAPI {
	return MockACMDeleteCertificateAPI(func(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error) {
		var availableCertificates map[string]bool = map[string]bool{}

		for _, mp := range mockParams {
			availableCertificates[mp.Arn] = true
		}

		if _, ok := availableCertificates[*params.CertificateArn]; !ok {
			return nil, errors.New("not exists")
		}

		return &acm.DeleteCertificateOutput{}, nil
	})
}

// GenerateMockACMRequestCertificateAPI returns MockACMRequestCertificateAPI
func GenerateMockACMRequestCertificateAPI(mockParams []MockParams) MockACMRequestCertificateAPI {
	return MockACMRequestCertificateAPI(func(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error) {
		return &acm.RequestCertificateOutput{}, nil
	})
}
