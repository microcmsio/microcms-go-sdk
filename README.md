# microCMS Go SDK

[microCMS](https://document.microcms.io/manual/api-request) Go SDK.

## Tutorial

See [official tutorial](https://document.microcms.io/tutorial/go/go-top).

## Installation

```sh
$ go get github.com/microcmsio/microcms-go-sdk
```

## Usage

### Import

```go
import "github.com/microcmsio/microcms-go-sdk"
```

### Create client object

```go
serviceDomain := "YOUR_DOMAIN" // YOUR_DOMAIN is the XXXX part of XXXX.microcms.io
apiKey := "YOUR_API_KEY"
client := microcms.New(serviceDomain, apiKey)
```

### Example content definition

```go
type YourContent struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Body        string    `json:"body,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	PublishedAt time.Time `json:"publishedAt,omitempty"`
	RevisedAt   time.Time `json:"revisedAt,omitempty"`
}

type YourContentList struct {
	Contents   []Content
	TotalCount int
	Limit      int
	Offset     int
}
```

### Get content list

```go
var list YourContentList
err := client.List(
	microcms.ListParams{
		Endpoint: "endpoint",
	},
	&list,
)
println(list.Contents[0].Title)
```

### Get content list with parameters

```go
var list YourContentList
err := client.List(
	microcms.ListParams{
		Endpoint: "endpoint",
		DraftKey: "abcd",
		Limit:    100,
		Offset:   1,
		Orders:   []string{"createdAt"},
		Q:        "Hello",
		Fields:   []string{"id", "title"},
		IDs:      []string{"foo"},
		Filters:  "publishedAt[greater_than]2021-01-01",
		Depth:    1,
	},
	&list,
)
println(list.Contents[0].Title)
```

### Get single content

```go
var content YourContent
err := client.Get(
	microcms.GetParams{
		Endpoint:  "endpoint",
		ContentID: "my-content-id",
	},
	&content,
)
println(content.Title)
```

### Get single content with parameters

```go
var content YourContent
err := client.Get(
	microcms.GetParams{
		Endpoint:  "endpoint",
		ContentID: "my-content-id",
		DraftKey:  "abcd",
		Fields:    []string{"id", "title"},
		Depth:     1,
	},
	&content,
)
println(content.Title)
```

### Get object form content

```go
var content YourContent
err := client.Get(
	microcms.GetParams{
		Endpoint: "endpoint",
	},
	&content,
)
println(content.Title)
```

### Create content

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint: "endpoint",
	Content:  YourContent{
		Title: "content",
		Body:  "Hello, content!",
	},
})
println(createResult.ID)
```

### Create content with specified ID

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Title: "my content",
		Body:  "Hello, my content!",
	},
})
println(createResult.ID)
```

### Create draft content

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint: "endpoint",
	Status:   microcms.StatusDraft,
	Content:  YourContent{
		Title: "draft content",
		Body:  "Hello, draft content!",
	},
})
println(createResult.ID)
```

### Update content

```go
updateResult, err := client.Update(microcms.UpdateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Body: "Hello, new content!",
	},
})
println(updateResult.ID)
```

### Update object form content

```go
updateResult, err := client.Update(microcms.UpdateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Body: "Hello, new content!",
	},
})
println(updateResult.ID)
```

### Delete content

```go
err := client.Delete(microcms.DeleteParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
})
```

### Error Handling

```go
data, err := client.Get(ctx, "endpoint", nil)
if err != nil {
    if httpErr, ok := err.(*sdk.HttpResponseError); ok {
        fmt.Printf("HTTP Status Code: %d\n", httpErr.Response.StatusCode)
        fmt.Printf("Error Message: %s\n", httpErr.ErrorMessage)
    }
    return
}
```
