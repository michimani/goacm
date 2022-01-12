package goacm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/michimani/goacm"
	"github.com/stretchr/testify/assert"
)

func Test_GetCertificate(t *testing.T) {
	ap := []goacm.MockACMParams{
		{
			Certificate: goacm.Certificate{
				Arn:        "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName: "test.example.com",
				Status:     string(types.CertificateStatusIssued),
				Type:       string(types.CertificateTypeAmazonIssued),
			},
		},
	}

	cases := []struct {
		name      string
		acmClient func(t *testing.T) goacm.MockACMAPI
		arn       string
		wantErr   bool
		expect    goacm.Certificate
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(ap)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			wantErr: false,
			expect: goacm.Certificate{
				Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName:    "test.example.com",
				Status:        string(types.CertificateStatusIssued),
				Type:          string(types.CertificateTypeAmazonIssued),
				FailureReason: "",
			},
		},
		{
			name: "notFound",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(ap)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-found-arn",
			wantErr: true,
			expect:  goacm.Certificate{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := goacm.GetCertificate(ctx, tt.acmClient(t), tt.arn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func Test_ListCertificateSummaries(t *testing.T) {
	mp := []goacm.MockACMParams{}
	expect := []types.CertificateSummary{}
	for i := 0; i < 3; i++ {
		mp = append(mp, goacm.MockACMParams{
			Certificate: goacm.Certificate{
				Arn:        fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1)),
				DomainName: fmt.Sprintf("test%d.example.com", (i + 1)),
			},
		})

		expect = append(expect, types.CertificateSummary{
			CertificateArn: aws.String(fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1))),
			DomainName:     aws.String(fmt.Sprintf("test%d.example.com", (i + 1))),
		})
	}

	cases := []struct {
		name      string
		acmClient func(t *testing.T) goacm.MockACMAPI
		wantErr   bool
		expect    []types.CertificateSummary
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(mp)
			},
			wantErr: false,
			expect:  expect,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := goacm.ListCertificateSummaries(ctx, tt.acmClient(t))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func Test_ListCertificates(t *testing.T) {
	mp := []goacm.MockACMParams{}
	expect := []goacm.Certificate{}
	for i := 0; i < 3; i++ {
		mp = append(mp, goacm.MockACMParams{
			Certificate: goacm.Certificate{
				Arn:              fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1)),
				DomainName:       fmt.Sprintf("test%d.example.com", (i + 1)),
				Status:           string(types.CertificateStatusIssued),
				Type:             string(types.CertificateTypeAmazonIssued),
				ValidationMethod: string(types.ValidationMethodDns),
				ValidationRecordSet: goacm.RecordSet{
					HostedDomainName: "example.com",
					Name:             fmt.Sprintf("_validation.%d.name.test.example.com", (i + 1)),
					Value:            fmt.Sprintf("_validation.%d.value.test.example.com", (i + 1)),
					Type:             string(route53Types.RRTypeCname),
				},
			},
		})

		expect = append(expect, goacm.Certificate{
			Arn:              fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1)),
			DomainName:       fmt.Sprintf("test%d.example.com", (i + 1)),
			Status:           string(types.CertificateStatusIssued),
			Type:             string(types.CertificateTypeAmazonIssued),
			ValidationMethod: string(types.ValidationMethodDns),
			ValidationRecordSet: goacm.RecordSet{
				HostedDomainName: "example.com",
				Name:             fmt.Sprintf("_validation.%d.name.test.example.com", (i + 1)),
				Value:            fmt.Sprintf("_validation.%d.value.test.example.com", (i + 1)),
				Type:             string(route53Types.RRTypeCname),
			},
		})
	}

	cases := []struct {
		name      string
		acmClient func(t *testing.T) goacm.MockACMAPI
		wantErr   bool
		expect    []goacm.Certificate
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(mp)
			},
			wantErr: false,
			expect:  expect,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := goacm.ListCertificates(ctx, tt.acmClient(t))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func Test_DeleteCertificate(t *testing.T) {
	ap := []goacm.MockACMParams{
		{
			Certificate: goacm.Certificate{
				Arn:              "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				Type:             string(types.CertificateTypeAmazonIssued),
				ValidationMethod: string(types.ValidationMethodDns),
				ValidationRecordSet: goacm.RecordSet{
					HostedDomainName: "example.com",
					Name:             "_validation.name.test.example.com",
					Value:            "_validation.value.test.example.com",
					Type:             string(route53Types.RRTypeCname),
				},
			},
		},
		{
			Certificate: goacm.Certificate{
				Arn:              "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				Type:             string(types.CertificateTypeAmazonIssued),
				ValidationMethod: string(types.ValidationMethodDns),
				ValidationRecordSet: goacm.RecordSet{
					HostedDomainName: "bothc-public-private.example.com",
					Name:             "_validation.name.test.example.com",
					Value:            "_validation.value.test.example.com",
					Type:             string(route53Types.RRTypeCname),
				},
			},
		},
	}

	rp := []goacm.MockRoute53Params{
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "example.com",
				Name:             "_validation.name.test.example.com",
				Value:            "_validation.value.test.example.com",
				Type:             string(route53Types.RRTypeCname),
			},
			ChangeAction: route53Types.ChangeActionDelete,
		},
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "bothc-public-private.example.com",
				Name:             "_validation.name.test.example.com",
				Value:            "_validation.value.test.example.com",
				Type:             string(route53Types.RRTypeCname),
			},
			ChangeAction: route53Types.ChangeActionDelete,
		},
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "bothc-public-private.example.com",
				Name:             "_validation.name.test.example.com",
				Value:            "_validation.value.test.example.com",
				Type:             string(route53Types.RRTypeCname),
			},
			ChangeAction:        route53Types.ChangeActionDelete,
			IsPrivateHostedZone: true,
		},
	}

	cases := []struct {
		name          string
		acmClient     func(t *testing.T) goacm.MockACMAPI
		route53Client func(t *testing.T) goacm.MockRoute53API
		arn           string
		wantErr       bool
		expect        *acm.DeleteCertificateOutput
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(ap)
			},
			route53Client: func(t *testing.T) goacm.MockRoute53API {
				return goacm.NewMockRoute53API(rp)
			},
			arn:     ap[0].Certificate.Arn,
			wantErr: false,
			expect:  &acm.DeleteCertificateOutput{},
		},
		{
			name: "normal: exists in both of public and private",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(ap)
			},
			route53Client: func(t *testing.T) goacm.MockRoute53API {
				return goacm.NewMockRoute53API(rp)
			},
			arn:     ap[1].Certificate.Arn,
			wantErr: false,
			expect:  &acm.DeleteCertificateOutput{},
		},
		{
			name: "notExists",
			acmClient: func(t *testing.T) goacm.MockACMAPI {
				return goacm.NewMockACMAPI(ap)
			},
			route53Client: func(t *testing.T) goacm.MockRoute53API {
				return goacm.NewMockRoute53API(rp)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-exists-arn",
			wantErr: true,
			expect:  nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			err := goacm.DeleteCertificate(ctx, tt.acmClient(t), tt.route53Client(t), tt.arn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_getPublicHostedZoneIDByDomainName(t *testing.T) {
	rp := []goacm.MockRoute53Params{
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "example.com",
			},
		},
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "bothc-public-private.example.com",
			},
		},
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "bothc-public-private.example.com",
			},
			IsPrivateHostedZone: true,
		},
		{
			RecordSet: goacm.RecordSet{
				HostedDomainName: "private.example.com",
			},
			IsPrivateHostedZone: true,
		},
	}

	route53API := func(t *testing.T) goacm.MockRoute53API {
		return goacm.NewMockRoute53API(rp)
	}

	cases := []struct {
		name       string
		domainName string
		wantErr    bool
		expect     string
	}{
		{
			name:       "normal: only public",
			domainName: rp[0].RecordSet.HostedDomainName,
			expect:     "example-com",
		},
		{
			name:       "normal: both of public and private",
			domainName: rp[1].RecordSet.HostedDomainName,
			expect:     "bothc-public-private-example-com",
		},
		{
			name:       "error: only private",
			domainName: rp[3].RecordSet.HostedDomainName,
			expect:     "",
		},
		{
			name:       "error: not exists",
			domainName: "not-exists.example.com",
			expect:     "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			id, err := goacm.ExportedGetPublicHostedZoneIDByDomainName(context.TODO(), route53API(tt), c.domainName)
			if c.wantErr {
				assert.Error(tt, err)
				assert.Empty(tt, id)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, c.expect, id)
		})
	}
}
