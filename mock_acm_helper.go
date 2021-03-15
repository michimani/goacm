package goacm

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// NewMockACMAPI return MockACMAPI.
func NewMockACMAPI(mockParams []MockACMParams) MockACMAPI {
	return MockACMAPI{
		DescribeCertificateAPI: NewMockACMDescribeCertificateAPI(mockParams),
		ListCertificatesAPI:    NewMockACMListCertificatesAPI(mockParams),
		DeleteCertificateAPI:   NewMockACMDeleteCertificateAPI(mockParams),
		RequestCertificateAPI:  NewMockACMRequestCertificateAPI(mockParams),
	}
}

// NewMockACMDescribeCertificateAPI returns MockACMDescribeCertificateAPI.
func NewMockACMDescribeCertificateAPI(mockParams []MockACMParams) MockACMDescribeCertificateAPI {
	return MockACMDescribeCertificateAPI(func(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
		if params.CertificateArn == nil {
			return nil, errors.New("expect Certificate ARN to not be nil")
		}

		var availableCertificates map[string]*acm.DescribeCertificateOutput = map[string]*acm.DescribeCertificateOutput{}
		for _, mp := range mockParams {
			if mp.Certificate.Arn == "" {
				continue
			}

			dv := types.DomainValidation{
				ValidationMethod: types.ValidationMethod(mp.Certificate.ValidationMethod),
			}
			if mp.Certificate.ValidationMethod == string(types.ValidationMethodDns) {
				dv.DomainName = aws.String(mp.Certificate.DomainName)
				dv.ValidationDomain = aws.String(mp.Certificate.ValidationRecordSet.HostedDomainName)
				dv.ResourceRecord = &types.ResourceRecord{
					Name:  aws.String(mp.Certificate.ValidationRecordSet.Name),
					Value: aws.String(mp.Certificate.ValidationRecordSet.Value),
					Type:  types.RecordType(mp.Certificate.ValidationRecordSet.Type),
				}
			}

			availableCertificates[mp.Certificate.Arn] = &acm.DescribeCertificateOutput{
				Certificate: &types.CertificateDetail{
					CertificateArn:          aws.String(mp.Certificate.Arn),
					DomainName:              aws.String(mp.Certificate.DomainName),
					Status:                  types.CertificateStatus(mp.Certificate.Status),
					Type:                    types.CertificateType(mp.Certificate.Type),
					FailureReason:           types.FailureReason(mp.Certificate.FailureReason),
					DomainValidationOptions: []types.DomainValidation{dv},
				},
			}
		}

		dco := availableCertificates[*params.CertificateArn]
		if dco == nil {
			return nil, fmt.Errorf("certificate arn not found arn: %s", *params.CertificateArn)
		}

		return dco, nil
	})
}

// NewMockACMListCertificatesAPI returns MockACMDescribeCertificateAPI.
func NewMockACMListCertificatesAPI(mockParams []MockACMParams) MockACMListCertificatesAPI {
	return MockACMListCertificatesAPI(func(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
		var csList []types.CertificateSummary

		for _, mp := range mockParams {
			csList = append(csList, types.CertificateSummary{
				CertificateArn: aws.String(mp.Certificate.Arn),
				DomainName:     aws.String(mp.Certificate.DomainName),
			})
		}

		return &acm.ListCertificatesOutput{
			CertificateSummaryList: csList,
		}, nil
	})
}

// NewMockACMDeleteCertificateAPI returns MockACMDeleteCertificateAPI
func NewMockACMDeleteCertificateAPI(mockParams []MockACMParams) MockACMDeleteCertificateAPI {
	return MockACMDeleteCertificateAPI(func(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error) {
		var availableCertificates map[string]bool = map[string]bool{}

		for _, mp := range mockParams {
			availableCertificates[mp.Certificate.Arn] = true
		}

		if _, ok := availableCertificates[*params.CertificateArn]; !ok {
			return nil, fmt.Errorf("certificate arn not found arn: %s", *params.CertificateArn)
		}

		return &acm.DeleteCertificateOutput{}, nil
	})
}

// NewMockACMRequestCertificateAPI returns MockACMRequestCertificateAPI
func NewMockACMRequestCertificateAPI(mockParams []MockACMParams) MockACMRequestCertificateAPI {
	return MockACMRequestCertificateAPI(func(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error) {
		return &acm.RequestCertificateOutput{}, nil
	})
}
