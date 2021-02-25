package gacm

import (
	"testing"

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
