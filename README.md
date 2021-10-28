# microCMS Golang SDK

It helps you to use microCMS from Golang applications.

## Getting Started

### Install

Install package.

```bash
$ go get github.com/microcmsio/microcms-go-sdk/v0
```

### How to use

#### Import

```go
import "github.com/microcmsio/microcms-go-sdk/v0"
```

#### Create client object

```go
serviceDomain := "YOUR_DOMAIN" // YOUR_DOMAIN is the XXXX part of XXXX.microcms.io
apiKey := "YOUR_API_KEY"
client := microcms.New(serviceDomain, apiKey)
```

#### Example content definition

```go
type YourContent struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Body        string     `json:"body,omitempty"`
	CreatedAt   time.Date  `json:"createdAt,omitempty"`
	UpdatedAt   time.Date  `json:"updatedAt,omitempty"`
	PublishedAt *time.Date `json:"publishedAt,omitempty"`
	RevisedAt   *time.Date `json:"revisedAt,omitempty"`
}

type YourContentList struct {
	Contents   []Content
	TotalCount int
	Limit      int
	Offset     int
}
```

#### Get content list

```go
var list YourContentList
err := client.List(
	microcms.ListParams{
		Endpoint: "endpoint",
		DraftKey: "abcd",                                 // Optional
		Limit:    100,                                    // Optional
		Offset:   1,                                      // Optional
		Orders:   []string{"createdAt"},                  // Optional
		Q:        "Hello",                                // Optional
		Fields:   []string{"id", "title"},                // Optional
		IDs:      []string{"foo"},                        // Optional
		Filters:  "publishedAt[greater_than]2021-01-01",  // Optional
		Depth:    1,                                      // Optional
	},
	&list,
)
println(list.Contents[0].Title)
```

#### Get single content

```go
var content YourContent
err := client.Get(
	microcms.GetParams{
		Endpoint:  "endpoint",
		ContentID: "my-content-id",
		DraftKey:  "abcd",                   // Optional
		Fields:    []string{"id", "title"},  // Optional
		Depth:     1,                        // Optional
	},
	&content,
)
println(content.Title)
```

#### Create content

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint:  "endpoint",
	ContentID: "draft-content-id",    // Optional
	Status:    microcms.StatusDraft,  // Optional
	Content: YourContent{
		Title: "draft-content",
		Body:  "Hello, draft content!",
	},
})
println(createResult.ID)
```

#### Update content

```go
updateResult, err := client.Update(microcms.UpdateParams{
	Endpoint:  "endpoint",
	ContentID: postResult.ID,
	Content: YourContent{
		Body: "Hello, new content!",
	},
})
println(updateResult.ID)
```

#### Delete

```go
err := client.Delete(microcms.DeleteParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
})
```
