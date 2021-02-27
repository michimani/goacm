package acmgo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCertificate(t *testing.T) {
	params := []MockParams{
		{
			Arn:             "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			DomainName:      "test.example.com",
			Status:          "ISSUED",
			CertificateType: "AMAZON_ISSUED",
		},
	}

	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMAPI
		arn     string
		wantErr bool
		expect  Certificate
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI(params)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			wantErr: false,
			expect: Certificate{
				Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
				DomainName:    "test.example.com",
				Status:        "ISSUED",
				Type:          "AMAZON_ISSUED",
				FailureReason: "",
			},
		},
		{
			name: "notFound",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI(params)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-found-arn",
			wantErr: true,
			expect:  Certificate{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := GetCertificate(tt.client(t), tt.arn)
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
	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMAPI
		wantErr bool
		expect  []types.CertificateSummary
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI([]MockParams{
					{
						Arn:        "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-1",
						DomainName: "test1.example.com",
					},
					{
						Arn:        "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-2",
						DomainName: "test2.example.com",
					},
					{
						Arn:        "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-3",
						DomainName: "test3.example.com",
					},
				})
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
		t.Run(tt.name, func(t *testing.T) {
			c, err := ListCertificateSummaries(tt.client(t))
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
	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMAPI
		wantErr bool
		expect  []Certificate
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI([]MockParams{
					{
						Arn:             "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-1",
						DomainName:      "test1.example.com",
						Status:          "ISSUED",
						CertificateType: "AMAZON_ISSUED",
					},
					{
						Arn:             "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-2",
						DomainName:      "test2.example.com",
						Status:          "ISSUED",
						CertificateType: "AMAZON_ISSUED",
					},
					{
						Arn:             "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-3",
						DomainName:      "test3.example.com",
						Status:          "ISSUED",
						CertificateType: "AMAZON_ISSUED",
					},
				})
			},
			wantErr: false,
			expect: []Certificate{
				{
					Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-1",
					DomainName:    "test1.example.com",
					Status:        "ISSUED",
					Type:          "AMAZON_ISSUED",
					FailureReason: "",
				},
				{
					Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-2",
					DomainName:    "test2.example.com",
					Status:        "ISSUED",
					Type:          "AMAZON_ISSUED",
					FailureReason: "",
				},
				{
					Arn:           "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn-3",
					DomainName:    "test3.example.com",
					Status:        "ISSUED",
					Type:          "AMAZON_ISSUED",
					FailureReason: "",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ListCertificates(tt.client(t))
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
	params := []MockParams{
		{
			Arn: "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
		},
	}

	cases := []struct {
		name    string
		client  func(t *testing.T) MockACMAPI
		arn     string
		wantErr bool
		expect  *acm.DeleteCertificateOutput
	}{
		{
			name: "normal",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI(params)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/this-is-a-sample-arn",
			wantErr: false,
			expect:  &acm.DeleteCertificateOutput{},
		},
		{
			name: "notExists",
			client: func(t *testing.T) MockACMAPI {
				return GenerateMockACMAPI(params)
			},
			arn:     "arn:aws:acm:ap-northeast-1:000000000000:certificate/not-exists-arn",
			wantErr: true,
			expect:  nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteCertificate(tt.client(t), tt.arn)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
