package acmgo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCertificates(t *testing.T) {
	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMDescribeCertificateAPI
		arn     string
		wantErr bool
		expect  Certificate
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMDescribeCertificateAPI {
				return GenerateMockACMDescribeCertificateAPI(
					t,
					"arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
					"test.example.com",
					"ISSUED",
					"AMAZON_ISSUED",
					"",
				)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			wantErr: false,
			expect: Certificate{
				ARN:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName:    "test.example.com",
				Status:        "ISSUED",
				Type:          "AMAZON_ISSUED",
				FailureReason: "",
			},
		},
		{
			name: "notFound",
			client: func(t *testing.T) MockACMDescribeCertificateAPI {
				return GenerateMockACMDescribeCertificateAPI(
					t,
					"arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
					"test.example.com",
					"ISSUED",
					"AMAZON_ISSUED",
					"",
				)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-found-arn",
			wantErr: true,
			expect:  Certificate{},
		},
	}

	for _, tt := range cases {
		c, err := GetCertificate(tt.client(t), tt.arn)
		if tt.wantErr {
			assert.Error(t, err)
			return
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.expect, c)
	}
}

func TestListCertificateSummaries(t *testing.T) {
	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMListCertificatesAPI
		wantErr bool
		expect  []types.CertificateSummary
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMListCertificatesAPI {
				return GenerateMockACMListCertificatesAPI(
					t,
					"arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-",
					"example.com",
					3,
				)
			},
			wantErr: false,
			expect: []types.CertificateSummary{
				{
					CertificateArn: aws.String("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-1"),
					DomainName:     aws.String("test1.example.com"),
				},
				{
					CertificateArn: aws.String("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-2"),
					DomainName:     aws.String("test2.example.com"),
				},
				{
					CertificateArn: aws.String("arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-3"),
					DomainName:     aws.String("test3.example.com"),
				},
			},
		},
	}

	for _, tt := range cases {
		c, err := ListCertificateSummaries(tt.client(t))
		if tt.wantErr {
			assert.Error(t, err)
			return
		}
		assert.NoError(t, err)
		assert.Equal(t, tt.expect, c)
	}
}
