package acmgo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// ACMAPI is an interface that defines ACM API.
type ACMAPI interface {
	ACMListCertificatesAPI
	ACMDescribeCertificateAPI
}

// ACMListCertificatesAPI is an interface that defines the set of ACM API operations required by the ListCertificates function.
type ACMListCertificatesAPI interface {
	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)
}

// ACMDescribeCertificateAPI is an interface that defines the set of ACM API operations required by the GetCertificate function.
type ACMDescribeCertificateAPI interface {
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)
}

// Certificate is a structure that represents a Certificate.
type Certificate struct {
	Arn           string
	Region        string
	DomainName    string
	Type          types.CertificateType
	Status        types.CertificateStatus
	FailureReason types.FailureReason
}
