acmgo
===

acmgo is a simple package for using AWS Certificate Manager from applications implimented Golang.

# Features

- List Certificates
- Get a Certificate

### TODO

- Publish Certificate

# Example

```go
package main

import (
	"fmt"

	"github.com/michimani/acmgo"
)

func main() {
	acmg, err := acmgo.NewACMgo("ap-northeast-1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	certificates, err := acmgo.ListCertificates(acmg.Client)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, c := range certificates {
		fmt.Println(c.Arn)
	}
	return
}
```

```bash
$ go run main.go

arn:aws:acm:ap-northeast-1:000000000000:certificate/00000000-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000001-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000002-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000003-xxxx-xxxx-0000-xxxxxxxxxxxx
arn:aws:acm:ap-northeast-1:000000000000:certificate/00000004-xxxx-xxxx-0000-xxxxxxxxxxxx
```