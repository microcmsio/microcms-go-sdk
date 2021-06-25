package main

import (
	"fmt"
	"github.com/microcmsio/microcms-go-sdk"
)

func main() {
	serviceDomain := "YOUR_DOMAIN"
	apiKey := "YOUR_API_KEY"

	c := microcms.CreateClient(serviceDomain, apiKey)
	data, _ := c.Get("endpoint")

	fmt.Println(data)
}
