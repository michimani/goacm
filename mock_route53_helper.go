package goacm

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// GenerateMockRoute53API returns MockRoute53API.
func GenerateMockRoute53API(mockParams []MockRoute53Params) MockRoute53API {
	return MockRoute53API{
		ListHostedZonesAPI:          GenerateMockListHostedZonesAPI(mockParams),
		ChangeResourceRecordSetsAPI: GenerateMockChangeResourceRecordSetsAPI(mockParams),
	}
}

// GenerateMockListHostedZonesAPI returns MockListHostedZonesAPI.
func GenerateMockListHostedZonesAPI(mockParams []MockRoute53Params) MockListHostedZonesAPI {
	return MockListHostedZonesAPI(func(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
		out := route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{},
		}

		for _, p := range mockParams {
			out.HostedZones = append(out.HostedZones, types.HostedZone{
				Id:   aws.String(strings.Replace(p.HostedDomainName, ".", "-", -1)),
				Name: aws.String(p.HostedDomainName + "."),
			})
		}

		return &out, nil
	})
}

// GenerateMockChangeResourceRecordSetsAPI returns MockChangeResourceRecordSetsAPI.
func GenerateMockChangeResourceRecordSetsAPI(mockParams []MockRoute53Params) MockChangeResourceRecordSetsAPI {
	return MockChangeResourceRecordSetsAPI(func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
		out := route53.ChangeResourceRecordSetsOutput{}

		return &out, nil
	})
}
