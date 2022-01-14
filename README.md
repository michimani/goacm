goacm
===

goacm is a simple package for using AWS Certificate Manager from applications implimented golang.

# Features

- List Certificates
- Get a Certificate
- Delete a Certificate
	- with Route 53 RecordSet that validates the domain (if validation method is DNS)
- Issue an SSL Certificate
	- Create Certificate
	- Create Route 53 RecordSet for validating the domain (if validation method is DNS)

# Example

## Create goacm client

```go
ctx := context.TODO()
g, err := goacm.NewGoACM(ctx, "ap-northeast-1")
if err != nil {
	fmt.Println(err.Error())
	return
}
```

## List Certificates

```go
ctx := context.TODO()
if certificates, err := goacm.ListCertificates(ctx, g.ACMClient); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("DomainName\tStatus\tARN")
	for _, c := range certificates {
		fmt.Printf("%s\t%s\t%s\n", c.DomainName, c.Status, c.Arn)
	}
}
```

## Get a Certificate

```go
arn := "arn:aws:acm:ap-northeast-1:000000000000:certificate/xxxxxxxx-1111-1111-1111-11111111xxxx"
ctx := context.TODO()
c, err := goacm.GetCertificate(ctx, g.ACMClient, arn)
if err != nil {
	fmt.Println(err.Error())
	return
}

fmt.Println("DomainName\tStatus\tARN")
fmt.Printf("%s\t%s\t%s\n", c.DomainName, c.Status, c.Arn)
```

## Issue a SSL Certificate

Request an ACM Certificate and create a RecordSet in Route 53 to validate the domain.

```go
method := "DNS"
targetDomain := "sample.exapmle.com"
hostedDomain := "example.com"
ctx := context.TODO()
res, err := goacm.IssueCertificate(ctx, g.ACMClient, g.Route53Client, method, targetDomain, hostedDomain)
if err != nil {
	fmt.Println(err.Error())
	return
}

fmt.Printf("ARN: %v", res.CertificateArn)
```

## Delete a Certificate

Delete the Route 53 RecordSet that was created for ACM Certificate and Domain validation.

```go
arn := "arn:aws:acm:ap-northeast-1:000000000000:certificate/xxxxxxxx-1111-1111-1111-11111111xxxx"
ctx := context.TODO()
if err := goacm.DeleteCertificate(ctx, g.ACMClient, g.Route53Client, arn); err != nil {
	fmt.Println(err.Error())
}
```
