package main

import (
	"fmt"
	"time"

	"github.com/microcmsio/microcms-go-sdk"
)

type ContentBase struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RevisedAt   *time.Time
	PublishedAt *time.Time
}
type ContentListBase struct {
	TotalCount int
	Offset     int
	Limit      int
}

type Blog struct {
	ContentBase
	Title string
	Body  string
}

type BlogList struct {
	ContentListBase
	Contents []Blog
}

func main() {
	serviceDomain := "YOUR_DOMAIN"
	apiKey := "YOUR_API_KEY"

	c := microcms.CreateClient(serviceDomain, apiKey)
	data := new(BlogList)
	_ = c.Get("endpoint", data)

	fmt.Printf("%+v\n", data)
}
