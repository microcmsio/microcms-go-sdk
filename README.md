# microCMS Golang SDK

It helps you to use microCMS from Golang applications.

## Getting Started

### Install

Install package.

```bash
$ go get github.com/microcmsio/microcms-go-sdk 
```

### How to use

First, create client.

```go
package main

import (
	"fmt"
	"github.com/microcmsio/microcms-go-sdk"
)

func main() {
	serviceDomain := "YOUR_DOMAIN" // YOUR_DOMAIN is the XXXX part of XXXX.microcms.io
	apiKey := "YOUR_API_KEY"
	globalDraftKey := "YOUR_GLOBAL_DRAFT_KEY" // If need 

	// If you specify globalDraftKey, please use microcms.GlobalDraftKey
	c := microcms.CreateClient(serviceDomain, apiKey, microcms.GlobalDraftKey(globalDraftKey))
}
```

After, How to use it below.

```go
	endpoint := "endpoint"
	contenttId := "contenttId" 

	data, _ := c.Get(endpoint, microcms.ContentId(contentId))
```
