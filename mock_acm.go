package goacm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
)

// MockACMParams is a structure with the elements needed to generate a mock.
type MockACMParams struct {
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
