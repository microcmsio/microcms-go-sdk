package main

import (
	"fmt"
	"github.com/wantainc/microcms-go-sdk"
)

func main() {
	serviceDomain := "YOUR_DOMAIN"
	apiKey := "YOUR_API_KEY"
	globalDraftKey := "YOUR_GLOBAL_DRAFT_KEY"

	c := microcms.CreateClient(serviceDomain, apiKey, globalDraftKey)
	data, _ := c.Get("endpoint")

	fmt.Println(data)
}
