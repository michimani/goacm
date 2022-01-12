package main

import (
	"context"
	"fmt"

	"github.com/michimani/goacm"
)

func main() {
	gacm, err := goacm.NewGoACM(context.TODO(), "ap-northeast-1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.TODO()
	listCertificate(ctx, gacm)
}

// List Certificate
func listCertificate(ctx context.Context, g *goacm.GoACM) {
	if certificates, err := goacm.ListCertificates(ctx, g.ACMClient); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DomainName\tStatus\tARN")
		for _, c := range certificates {
			fmt.Printf("%s\t%s\t%s\n", c.DomainName, c.Status, c.Arn)
		}
	}
}

// Get a Certificate
func getCertificate(ctx context.Context, g *goacm.GoACM, arn string) {
	c, err := goacm.GetCertificate(ctx, g.ACMClient, arn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("DomainName\tStatus\tARN")
	fmt.Printf("%s\t%s\t%s\n", c.DomainName, c.Status, c.Arn)
}

// Issue a Certificate
func issueCertificate(ctx context.Context, g *goacm.GoACM, targetDomain, hostedDomain, method string) {
	res, err := goacm.IssueCertificate(ctx, g.ACMClient, g.Route53Client, method, targetDomain, hostedDomain)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("ARN: %v", res.CertificateArn)
}

// Delete a Certificate
func deleteCertificate(ctx context.Context, g *goacm.GoACM, arn string) {
	if err := goacm.DeleteCertificate(ctx, g.ACMClient, g.Route53Client, arn); err != nil {
		fmt.Println(err.Error())
	}
}
