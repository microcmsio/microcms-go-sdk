package main

import (
	"log"
	"os"
	"time"

	"github.com/microcmsio/microcms-go-sdk"
)

type Blog struct {
	ID          *string    `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	RevisedAt   *time.Time `json:"revisedAt,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Title       string     `json:"title,omitempty"`
	Contents    string     `json:"contents,omitempty"`
}

type BlogList struct {
	Contents   []Blog `json:"contents"`
	TotalCount int    `json:"totalCount"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

func main() {
	serviceDomain := os.Getenv("YOUR_DOMAIN")
	apiKey := os.Getenv("YOUR_API_KEY")

	client := microcms.New(serviceDomain, apiKey)

	var blogList BlogList
	if err := client.List(
		microcms.ListParams{
			Endpoint: "blog",
			Limit:    100,
			Offset:   1,
			Orders:   []string{"updatedAt"},
			Q:        "Hello",
			Fields:   []string{"id", "title"},
			Filters:  "publishedAt[greater_than]2021-01-01",
		},
		&blogList,
	); err != nil {
		panic(err)
	}
	log.Printf("%+v", blogList)

	var blog Blog
	if err := client.Get(
		microcms.GetParams{
			Endpoint:  "blog",
			ContentID: "my-content-id",
		},
		&blog,
	); err != nil {
		panic(err)
	}
	log.Printf("%+v", blog)

	postResult, err := client.Create(microcms.CreateParams{
		Endpoint: "blog",
		Content: Blog{
			Title:    "my-content",
			Contents: "Hello, POST request!",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", postResult)

	putResult, err := client.Create(
		microcms.CreateParams{
			Endpoint:  "blog",
			ContentID: "my-content-id",
			Content: Blog{
				Title:    "my-content",
				Contents: "Hello, PUT request!",
			},
		})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", putResult)

	draftResult, err := client.Create(microcms.CreateParams{
		Endpoint:  "blog",
		ContentID: "draft-content-id",
		Status:    microcms.StatusDraft,
		Content: Blog{
			Title:    "draft-content",
			Contents: "Hello, draft content!",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", draftResult)

	updateResult, err := client.Update(microcms.UpdateParams{
		Endpoint:  "blog",
		ContentID: postResult.ID,
		Content: Blog{
			Contents: "Hello, new content!",
		},
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", updateResult)

	if err := client.Delete(microcms.DeleteParams{
		Endpoint:  "blog",
		ContentID: "my-content-id",
	}); err != nil {
		panic(err)
	}
}
