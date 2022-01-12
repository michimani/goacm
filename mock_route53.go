package goacm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// MockRoute53Params is a structure with the elements needed to generate a mock.
type MockRoute53Params struct {
	RecordSet           RecordSet
	ChangeAction        types.ChangeAction
	IsPrivateHostedZone bool
}

// MockRoute53API is a struct that represents a Route 53 client.
type MockRoute53API struct {
	ListHostedZonesAPI          MockListHostedZonesAPI
	ListResourceRecordSetsAPI   MockListResourceRecordSetsAPI
	ChangeResourceRecordSetsAPI MockChangeResourceRecordSetsAPI
}

// MockListHostedZonesAPI is a type that represents a function that mock Route 53's MockListHostedZones.
type MockListHostedZonesAPI func(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)

// MockListResourceRecordSetsAPI is a type that represents a function that mock Route 53's MockListResourceRecordSets.
type MockListResourceRecordSetsAPI func(ctx context.Context, params *route53.ListResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error)

// MockChangeResourceRecordSetsAPI is a type that represents a function that mock Route 53's MockChangeResourceRecordSets.
type MockChangeResourceRecordSetsAPI func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)

// ListHostedZones returns a function that mock original of Route 53 ListHostedZones.
func (m MockRoute53API) ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
	return m.ListHostedZonesAPI(ctx, params, optFns...)
}

// ListResourceRecordSets returns a function that mock original of Route 53 ListResourceRecordSets.
func (m MockRoute53API) ListResourceRecordSets(ctx context.Context, params *route53.ListResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error) {
	return m.ListResourceRecordSetsAPI(ctx, params, optFns...)
}

// ChangeResourceRecordSets returns a function that mock original of Route 53 ChangeResourceRecordSets.
func (m MockRoute53API) ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	return m.ChangeResourceRecordSetsAPI(ctx, params, optFns...)
}
