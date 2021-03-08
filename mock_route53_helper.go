package goacm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// GenerateMockRoute53API returns MockRoute53API.
func GenerateMockRoute53API(mockParams []MockRoute53Params) MockRoute53API {
	return MockRoute53API{
		ListHostedZonesAPI:          GenerateMockListHostedZonesAPI(mockParams),
		ListResourceRecordSetsAPI:   GenerateMockListResourceRecordSetsAPI(mockParams),
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
				Id:   aws.String(strings.Replace(p.RecordSet.HostedDomainName, ".", "-", -1)),
				Name: aws.String(p.RecordSet.HostedDomainName + "."),
			})
		}

		return &out, nil
	})
}

// GenerateMockListResourceRecordSetsAPI returns MockListResourceRecordSetsAPI.
func GenerateMockListResourceRecordSetsAPI(mockParams []MockRoute53Params) MockListResourceRecordSetsAPI {
	return MockListResourceRecordSetsAPI(func(ctx context.Context, params *route53.ListResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error) {
		out := route53.ListResourceRecordSetsOutput{}

		available := map[string]*types.ResourceRecordSet{}
		for _, p := range mockParams {
			available[p.RecordSet.Name] = &types.ResourceRecordSet{
				Name: aws.String(p.RecordSet.Name),
				TTL:  aws.Int64(p.RecordSet.TTL),
				Type: types.RRTypeCname,
				ResourceRecords: []types.ResourceRecord{
					{
						Value: aws.String(p.RecordSet.Value),
					},
				},
			}
		}

		rrs := available[aws.ToString(params.StartRecordName)]
		if rrs == nil {
			return nil, fmt.Errorf("resource record sets not exists name: %s", *params.StartRecordName)
		}

		out.ResourceRecordSets = append(out.ResourceRecordSets, *rrs)
		return &out, nil
	})
}

// GenerateMockChangeResourceRecordSetsAPI returns MockChangeResourceRecordSetsAPI.
func GenerateMockChangeResourceRecordSetsAPI(mockParams []MockRoute53Params) MockChangeResourceRecordSetsAPI {
	return MockChangeResourceRecordSetsAPI(func(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
		out := route53.ChangeResourceRecordSetsOutput{}

		available := map[string]*types.ChangeBatch{}
		for _, p := range mockParams {
			hostedZoneID := strings.Replace(p.RecordSet.HostedDomainName, ".", "-", -1)
			available[hostedZoneID] = &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: p.ChangeAction,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(p.RecordSet.Name),
							Type: types.RRType(p.RecordSet.Type),
							ResourceRecords: []types.ResourceRecord{
								{
									Value: aws.String(p.RecordSet.Value),
								},
							},
						},
					},
				},
			}
		}

		cb := available[*params.HostedZoneId]
		if cb == nil {
			return nil, fmt.Errorf("hosted zone id not exists hostedzoneid: %s", *params.HostedZoneId)
		}

		if cb.Changes[0].Action != params.ChangeBatch.Changes[0].Action {
			return nil, fmt.Errorf("change action not matches: expected %v but %v", cb.Changes[0].Action, params.ChangeBatch.Changes[0].Action)
		}

		if *cb.Changes[0].ResourceRecordSet.Name != *params.ChangeBatch.Changes[0].ResourceRecordSet.Name {
			return nil, fmt.Errorf("record set name not mathces: expected %v but %v",
				*cb.Changes[0].ResourceRecordSet.Name,
				*params.ChangeBatch.Changes[0].ResourceRecordSet.Name)
		}

		if cb.Changes[0].ResourceRecordSet.Type != params.ChangeBatch.Changes[0].ResourceRecordSet.Type {
			return nil, fmt.Errorf("record set type not mathces: expected %v but %v",
				cb.Changes[0].ResourceRecordSet.Type,
				params.ChangeBatch.Changes[0].ResourceRecordSet.Type)
		}

		if *cb.Changes[0].ResourceRecordSet.ResourceRecords[0].Value != *params.ChangeBatch.Changes[0].ResourceRecordSet.ResourceRecords[0].Value {
			return nil, fmt.Errorf("record set resource record value not mathces: expected %v but %v",
				*cb.Changes[0].ResourceRecordSet.ResourceRecords[0].Value,
				*params.ChangeBatch.Changes[0].ResourceRecordSet.ResourceRecords[0].Value)
		}

		return &out, nil
	})
}
