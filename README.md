# microCMS Golang SDK

It helps you to use microCMS from Golang applications.

## Getting Started

### Install

Install package.

```bash
$ go get github.com/microcmsio/microcms-go-sdk 
```

### How to use

```go
package main

import (
	"fmt"
	"github.com/microcmsio/microcms-go-sdk"
)

type Content struct {
	ID          string
	Title       string
	Body        string
	CreatedAt   time.Date
	UpdatedAt   time.Date
	PublishedAt *time.Date
	RevisedAt   *time.Date
}

func main() {
	serviceDomain := "YOUR_DOMAIN" // YOUR_DOMAIN is the XXXX part of XXXX.microcms.io
	apiKey := "YOUR_API_KEY"
	globalDraftKey := "YOUR_GLOBAL_DRAFT_KEY" // If need 

	// First, create client.

	// If you specify globalDraftKey, please use microcms.GlobalDraftKey
	c := microcms.CreateClient(serviceDomain, apiKey, microcms.GlobalDraftKey(globalDraftKey))

	// After, How to use it below.

	endpoint := "endpoint"
	contenttId := "contenttId"
	data := new(Content)

	_ = c.Get(endpoint, data, microcms.ContentID(contentId))

	fmt.Printf("%+v\n", response)
}
```
