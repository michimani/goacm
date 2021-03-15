package goacm

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCertificate(t *testing.T) {
	ap := []MockACMParams{
		{
			Certificate: Certificate{
				Arn:        "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName: "test.example.com",
				Status:     string(types.CertificateStatusIssued),
				Type:       string(types.CertificateTypeAmazonIssued),
			},
		},
	}

	cases := []struct {
		name      string
		acmClient func(t *testing.T) MockACMAPI
		arn       string
		wantErr   bool
		expect    Certificate
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(ap)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			wantErr: false,
			expect: Certificate{
				Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName:    "test.example.com",
				Status:        string(types.CertificateStatusIssued),
				Type:          string(types.CertificateTypeAmazonIssued),
				FailureReason: "",
			},
		},
		{
			name: "notFound",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(ap)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-found-arn",
			wantErr: true,
			expect:  Certificate{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := GetCertificate(ctx, tt.acmClient(t), tt.arn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func TestListCertificateSummaries(t *testing.T) {
	mp := []MockACMParams{}
	expect := []types.CertificateSummary{}
	for i := 0; i < 3; i++ {
		mp = append(mp, MockACMParams{
			Certificate: Certificate{
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
		acmClient func(t *testing.T) MockACMAPI
		wantErr   bool
		expect    []types.CertificateSummary
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(mp)
			},
			wantErr: false,
			expect:  expect,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := ListCertificateSummaries(ctx, tt.acmClient(t))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func TestListCertificates(t *testing.T) {
	mp := []MockACMParams{}
	expect := []Certificate{}
	for i := 0; i < 3; i++ {
		mp = append(mp, MockACMParams{
			Certificate: Certificate{
				Arn:              fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1)),
				DomainName:       fmt.Sprintf("test%d.example.com", (i + 1)),
				Status:           string(types.CertificateStatusIssued),
				Type:             string(types.CertificateTypeAmazonIssued),
				ValidationMethod: string(types.ValidationMethodDns),
				ValidationRecordSet: RecordSet{
					HostedDomainName: "example.com",
					Name:             fmt.Sprintf("_validation.%d.name.test.example.com", (i + 1)),
					Value:            fmt.Sprintf("_validation.%d.value.test.example.com", (i + 1)),
					Type:             string(route53Types.RRTypeCname),
				},
			},
		})

		expect = append(expect, Certificate{
			Arn:              fmt.Sprintf("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-%d", (i + 1)),
			DomainName:       fmt.Sprintf("test%d.example.com", (i + 1)),
			Status:           string(types.CertificateStatusIssued),
			Type:             string(types.CertificateTypeAmazonIssued),
			ValidationMethod: string(types.ValidationMethodDns),
			ValidationRecordSet: RecordSet{
				HostedDomainName: "example.com",
				Name:             fmt.Sprintf("_validation.%d.name.test.example.com", (i + 1)),
				Value:            fmt.Sprintf("_validation.%d.value.test.example.com", (i + 1)),
				Type:             string(route53Types.RRTypeCname),
			},
		})
	}

	cases := []struct {
		name      string
		acmClient func(t *testing.T) MockACMAPI
		wantErr   bool
		expect    []Certificate
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(mp)
			},
			wantErr: false,
			expect:  expect,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			c, err := ListCertificates(ctx, tt.acmClient(t))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, c)
		})
	}
}

func TestDeleteCertificate(t *testing.T) {
	ap := []MockACMParams{
		{
			Certificate: Certificate{
				Arn:              "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				Type:             string(types.CertificateTypeAmazonIssued),
				ValidationMethod: string(types.ValidationMethodDns),
				ValidationRecordSet: RecordSet{
					HostedDomainName: "example.com",
					Name:             "_validation.name.test.example.com",
					Value:            "_validation.value.test.example.com",
					Type:             string(route53Types.RRTypeCname),
				},
			},
		},
	}

	rp := []MockRoute53Params{
		{
			RecordSet: RecordSet{
				HostedDomainName: "example.com",
				Name:             "_validation.name.test.example.com",
				Value:            "_validation.value.test.example.com",
				Type:             string(route53Types.RRTypeCname),
			},
			ChangeAction: route53Types.ChangeActionDelete,
		},
	}

	cases := []struct {
		name          string
		acmClient     func(t *testing.T) MockACMAPI
		route53Client func(t *testing.T) MockRoute53API
		arn           string
		wantErr       bool
		expect        *acm.DeleteCertificateOutput
	}{
		{
			name: "normal",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(ap)
			},
			route53Client: func(t *testing.T) MockRoute53API {
				return NewMockRoute53API(rp)
			},
			arn:     ap[0].Certificate.Arn,
			wantErr: false,
			expect:  &acm.DeleteCertificateOutput{},
		},
		{
			name: "notExists",
			acmClient: func(t *testing.T) MockACMAPI {
				return NewMockACMAPI(ap)
			},
			route53Client: func(t *testing.T) MockRoute53API {
				return NewMockRoute53API(rp)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-exists-arn",
			wantErr: true,
			expect:  nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			err := DeleteCertificate(ctx, tt.acmClient(t), tt.route53Client(t), tt.arn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
