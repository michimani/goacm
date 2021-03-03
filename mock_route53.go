package goacm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// MockRoute53Params is a structure with the elements needed to generate a mock.
type MockRoute53Params struct {
	HostedDomainName string
	RecordSet        RecordSet
}

// MockRoute53API is a struct that represents a Route 53 client.
type MockRoute53API struct {
	ListHostedZonesAPI          MockListHostedZonesAPI
	ChangeResourceRecordSetsAPI MockChangeResourceRecordSetsAPI
}

// MockListHostedZonesAPI is a type that represents a function that mock Route 53's MockListHostedZones.
type MockListHostedZonesAPI func(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)

// MockChangeResourceRecordSetsAPI is a type that represents a function that mock Route 53's MockChangeResourceRecordSets.
type MockChangeResourceRecordSetsAPI func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)

// ListHostedZones returns a function that mock original of Route 53 ListHostedZones.
func (m MockRoute53API) ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
	return m.ListHostedZonesAPI(ctx, params, optFns...)
}

// ChangeResourceRecordSets returns a function that mock original of Route 53 ChangeResourceRecordSets.
func (m MockRoute53API) ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	return m.ChangeResourceRecordSetsAPI(ctx, params, optFns...)
}
