package microcms

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_List(t *testing.T) {
	type Blog struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	type BlogList struct {
		Contents   []Blog
		TotalCount int
		Limit      int
		Offset     int
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodGet,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		resBody: `
			{
				"contents": [
					{
						"id": "foo",
						"title": "Hello, microCMS!",
						"createdAt": "2021-10-28T04:04:29.625Z",
						"updatedAt": "2021-10-28T04:04:29.625Z",
						"publishedAt": "2021-10-28T04:04:29.625Z",
						"revisedAt": "2021-10-28T04:04:29.625Z"
					}
				],
				"totalCount": 1,
				"limit": 10,
				"offset": 0
			}
		`,
	})

	blog := BlogList{}
	err := client.List(
		ListParams{
			Endpoint: "blog",
		},
		&blog,
	)
	require.NoError(t, err)
	require.EqualValues(
		t,
		BlogList{
			Contents:   []Blog{{ID: "foo", Title: "Hello, microCMS!"}},
			TotalCount: 1,
			Limit:      10,
			Offset:     0,
		},
		blog,
	)
}

func TestClient_List_querystring(t *testing.T) {
	type Blog struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	type BlogList struct {
		Contents   []Blog
		TotalCount int
		Limit      int
		Offset     int
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodGet,
		url: &url.URL{
			Scheme:   "https",
			Host:     "serviceDomain.microcms.io",
			Path:     "/api/v1/blog",
			RawQuery: "depth=1&draftKey=abcd&fields=id%2Ctitle&filters=publishedAt%5Bgreater_than%5D2021-01-01&ids=foo&limit=100&offset=1&orders=createdAt&q=Hello",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		resBody: `
			{
				"contents": [
					{
						"id": "foo",
						"title": "Hello, microCMS!",
						"createdAt": "2021-10-28T04:04:29.625Z",
						"updatedAt": "2021-10-28T04:04:29.625Z",
						"publishedAt": "2021-10-28T04:04:29.625Z",
						"revisedAt": "2021-10-28T04:04:29.625Z"
					}
				],
				"totalCount": 1,
				"limit": 10,
				"offset": 0
			}
		`,
	})

	blog := BlogList{}
	err := client.List(
		ListParams{
			Endpoint: "blog",
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
		&blog,
	)
	require.NoError(t, err)
	require.EqualValues(
		t,
		BlogList{
			Contents:   []Blog{{ID: "foo", Title: "Hello, microCMS!"}},
			TotalCount: 1,
			Limit:      10,
			Offset:     0,
		},
		blog,
	)
}

func TestClient_Get(t *testing.T) {
	type Blog struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodGet,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog/foo",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		resBody: `
			{
				"id": "foo",
				"title": "Hello, microCMS!",
				"createdAt": "2021-10-28T04:04:29.625Z",
				"updatedAt": "2021-10-28T04:04:29.625Z",
				"publishedAt": "2021-10-28T04:04:29.625Z",
				"revisedAt": "2021-10-28T04:04:29.625Z"
			}
		`,
	})

	blog := Blog{}
	err := client.Get(
		GetParams{
			Endpoint:  "blog",
			ContentID: "foo",
		},
		&blog,
	)
	require.NoError(t, err)
	require.EqualValues(
		t, Blog{
			ID:    "foo",
			Title: "Hello, microCMS!",
		},
		blog,
	)
}

func TestClient_Get_querystring(t *testing.T) {
	type Blog struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodGet,
		url: &url.URL{
			Scheme:   "https",
			Host:     "serviceDomain.microcms.io",
			Path:     "/api/v1/blog/foo",
			RawQuery: "depth=1&draftKey=abcd&fields=id%2Ctitle",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		resBody: `
			{
				"id": "foo",
				"title": "Hello, microCMS!",
				"createdAt": "2021-10-28T04:04:29.625Z",
				"updatedAt": "2021-10-28T04:04:29.625Z",
				"publishedAt": "2021-10-28T04:04:29.625Z",
				"revisedAt": "2021-10-28T04:04:29.625Z"
			}
		`,
	})

	blog := Blog{}
	err := client.Get(
		GetParams{
			Endpoint:  "blog",
			ContentID: "foo",
			DraftKey:  "abcd",
			Fields:    []string{"id", "title"},
			Depth:     1,
		},
		&blog,
	)
	require.NoError(t, err)
	require.EqualValues(
		t, Blog{
			ID:    "foo",
			Title: "Hello, microCMS!",
		},
		blog,
	)
}
