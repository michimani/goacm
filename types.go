package goacm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// ACMAPI is an interface that defines ACM API.
type ACMAPI interface {
	ACMListCertificatesAPI
	ACMDescribeCertificateAPI
	ACMDeleteCertificateAPI
	ACMRequestCertificateAPI
}

// Route53API is an interface that defines Route53 API.
type Route53API interface {
	Route53ListHostedZonesAPI
	Route53ChangeResourceRecordSetsAPI
}

// ACMListCertificatesAPI is an interface that defines the set of ACM API operations required by the ListCertificates function.
type ACMListCertificatesAPI interface {
	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)
}

// ACMDescribeCertificateAPI is an interface that defines the set of ACM API operations required by the DescribeCertificate function.
type ACMDescribeCertificateAPI interface {
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)
}

// ACMDeleteCertificateAPI is an interface that defines the set of ACM API operations required by the DeleteCertificate function.
type ACMDeleteCertificateAPI interface {
	DeleteCertificate(ctx context.Context, params *acm.DeleteCertificateInput, optFns ...func(*acm.Options)) (*acm.DeleteCertificateOutput, error)
}

// ACMRequestCertificateAPI is an interface that defines the set of ACM API operations required by the DeleteCertificate function.
type ACMRequestCertificateAPI interface {
	RequestCertificate(ctx context.Context, params *acm.RequestCertificateInput, optFns ...func(*acm.Options)) (*acm.RequestCertificateOutput, error)
}

// Route53ListHostedZonesAPI is an interface that defines the set of Route53 API operations required by the ListHostedZone function.
type Route53ListHostedZonesAPI interface {
	ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)
}

// Route53ChangeResourceRecordSetsAPI is an interface that defines the set of Route53 API operations required by the ChangeResourceRecordSets function.
type Route53ChangeResourceRecordSetsAPI interface {
	ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
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

// ValidationMethod is a type that represents the method when requesting a certificate.
type ValidationMethod types.ValidationMethod

const (
	// ValidationMethodDNS is a construct valur for validating certificate with DNS.
	ValidationMethodDNS ValidationMethod = "DNS"

	// ValidationMethodEmail is a construct valur for validating certificate with EMAIL.
	ValidationMethodEmail ValidationMethod = "EMAIL"
)

// IssueCertificateResult is a structure that represents a reault of IssueCertificate.
type IssueCertificateResult struct {
	CertificateArn        string
	DomainName            string
	HostedDomainName      string
	HosteZoneID           string
	ValidationMethod      string
	ValidationRecordName  string
	ValidationRecordValue string
}
