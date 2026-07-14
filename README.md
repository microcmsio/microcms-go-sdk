# microCMS Go SDK

[microCMS](https://document.microcms.io/manual/api-request) のGo SDKです。

## 保守方針

このSDKの現在の保守レベルは `Maintenance` です。

詳細は[SDKの保守方針](https://document.microcms.io/manual/limitations#hc2b0bc6659)をご覧ください。

## チュートリアル

[公式チュートリアル](https://document.microcms.io/tutorial/go/go-top)をご覧ください。

## インストール

```sh
$ go get github.com/microcmsio/microcms-go-sdk
```

## 使い方

### インポート

```go
import "github.com/microcmsio/microcms-go-sdk"
```

### クライアントオブジェクトの作成

```go
serviceDomain := "YOUR_DOMAIN" // YOUR_DOMAINはXXXX.microcms.ioのXXXXの部分です
apiKey := "YOUR_API_KEY"
client := microcms.New(serviceDomain, apiKey)
```

### コンテンツ定義の例

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

### コンテンツ一覧の取得

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

### パラメータを指定したコンテンツ一覧の取得

```go
var list YourContentList
err := client.List(
	microcms.ListParams{
		Endpoint: "endpoint",
		DraftKey: "abcd",
		Limit:    100,
		Offset:   1,
		Orders:   []string{"createdAt"},
		Q:        "こんにちは",
		Fields:   []string{"id", "title"},
		IDs:      []string{"foo"},
		Filters:  "publishedAt[greater_than]2021-01-01",
		Depth:    1,
	},
	&list,
)
println(list.Contents[0].Title)
```

### 単一コンテンツの取得

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

### パラメータを指定した単一コンテンツの取得

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

### オブジェクト形式のコンテンツの取得

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

### コンテンツの作成

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint: "endpoint",
	Content:  YourContent{
		Title: "コンテンツ",
		Body:  "こんにちは、コンテンツ！",
	},
})
println(createResult.ID)
```

### IDを指定したコンテンツの作成

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Title: "マイコンテンツ",
		Body:  "こんにちは、マイコンテンツ！",
	},
})
println(createResult.ID)
```

### 下書きコンテンツの作成

```go
createResult, err := client.Create(microcms.CreateParams{
	Endpoint: "endpoint",
	Status:   microcms.StatusDraft,
	Content:  YourContent{
		Title: "下書きコンテンツ",
		Body:  "こんにちは、下書きコンテンツ！",
	},
})
println(createResult.ID)
```

### コンテンツの更新

```go
updateResult, err := client.Update(microcms.UpdateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Body: "こんにちは、新しいコンテンツ！",
	},
})
println(updateResult.ID)
```

### オブジェクト形式のコンテンツの更新

```go
updateResult, err := client.Update(microcms.UpdateParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
	Content:   YourContent{
		Body: "こんにちは、新しいコンテンツ！",
	},
})
println(updateResult.ID)
```

### コンテンツの削除

```go
err := client.Delete(microcms.DeleteParams{
	Endpoint:  "endpoint",
	ContentID: "my-content-id",
})
```

### エラーハンドリング

```go
data, err := client.Get(ctx, "endpoint", nil)
if err != nil {
    if httpErr, ok := err.(*sdk.HttpResponseError); ok {
        fmt.Printf("HTTPステータスコード: %d\n", httpErr.Response.StatusCode)
        fmt.Printf("エラーメッセージ: %s\n", httpErr.ErrorMessage)
    }
    return
}
```
