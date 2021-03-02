goacm
===

goacm is a simple package for using AWS Certificate Manager from applications implimented golang.

# Features

- List Certificates
- Get a Certificate
- Delete a Certificate
- Issue an SSL Certificate

### TODO

- Delete a certificate with Route 53 Record sets that validate domain.

# Example

```go
package main

import (
	"fmt"

	"github.com/michimani/goacm"
)

func main() {
	gacm, err := goacm.NewGoACM("ap-northeast-1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// List certificates and print each ARN.
	fmt.Println("List certificates and print each ARN.")
	listCertificate(gacm)

	// Issue an SSL certificate.
	fmt.Println("Issue an SSL certificate.")
	issueCertificate(gacm)
}

func listCertificate(g *goacm.GoACM) {
	if certificates, err := goacm.ListCertificates(g.ACMClient); err != nil {
		fmt.Println(err.Error())
	} else {
		for _, c := range certificates {
			fmt.Println(c.Arn)
		}
	}
}

func issueCertificate(g *goacm.GoACM) {
	targetDomain := "goacm.example.com"
	hostedDomain := "example.com"
	var validationMethod goacm.ValidationMethod = "DNS"
	if res, err := goacm.IssueCertificate(g.ACMClient, g.Route53Client, validationMethod, targetDomain, hostedDomain); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%v", res)
	}
}
```

```bash
$ go run main.go

List certificates and print each ARN.
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000000-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000001-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000002-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000003-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000004-xxxx-xxxx-0000-xxxxxxxxxxxx
Issue an SSL certificate.
{arn:aws:acm:ap-northeast-1:000000000000:certificate/00000005-xxxx-xxxx-0000-xxxxxxxxxxxx goacm.example.com example.com /hostedzone/Z3XXXXXXXXXXXX DNS _32xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.goacm.example.com. _80xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.xxxxxxxxxx.acm-validations.aws.}
```