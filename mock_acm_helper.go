package goacm

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// GenerateMockACMAPI return MockACMAPI.
func GenerateMockACMAPI(mockParams []MockACMParams) MockACMAPI {
	return MockACMAPI{
		DescribeCertificateAPI: GenerateMockACMDescribeCertificateAPI(mockParams),
		ListCertificatesAPI:    GenerateMockACMListCertificatesAPI(mockParams),
		DeleteCertificateAPI:   GenerateMockACMDeleteCertificateAPI(mockParams),
		RequestCertificateAPI:  GenerateMockACMRequestCertificateAPI(mockParams),
	}
}

// GenerateMockACMDescribeCertificateAPI returns MockACMDescribeCertificateAPI.
func GenerateMockACMDescribeCertificateAPI(mockParams []MockACMParams) MockACMDescribeCertificateAPI {
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
func GenerateMockACMListCertificatesAPI(mockParams []MockACMParams) MockACMListCertificatesAPI {
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
func GenerateMockACMDeleteCertificateAPI(mockParams []MockACMParams) MockACMDeleteCertificateAPI {
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
func GenerateMockACMRequestCertificateAPI(mockParams []MockACMParams) MockACMRequestCertificateAPI {
	return MockACMRequestCertificateAPI(func(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error) {
		return &acm.RequestCertificateOutput{}, nil
	})
}
