package microcms

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Create_Post(t *testing.T) {
	type Blog struct {
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodPost,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		reqBody:    `{"title": "Hello, microCMS!"}`,
		resBody:    `{"id": "foo"}`,
	})

	createResult, err := client.Create(CreateParams{
		Endpoint: "blog",
		Content: Blog{
			Title: "Hello, microCMS!",
		},
	})
	require.NoError(t, err)
	require.EqualValues(
		t,
		&CreateResponse{ID: "foo"},
		createResult,
	)
}

func TestClient_Create_Put(t *testing.T) {
	type Blog struct {
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodPut,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog/foo",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		reqBody:    `{"title": "Hello, microCMS!"}`,
		resBody:    `{"id": "foo"}`,
	})

	createResult, err := client.Create(CreateParams{
		Endpoint:  "blog",
		ContentID: "foo",
		Content: Blog{
			Title: "Hello, microCMS!",
		},
	})
	require.NoError(t, err)
	require.EqualValues(
		t,
		&CreateResponse{ID: "foo"},
		createResult,
	)
}

func TestClient_Create_Draft(t *testing.T) {
	type Blog struct {
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodPost,
		url: &url.URL{
			Scheme:   "https",
			Host:     "serviceDomain.microcms.io",
			Path:     "/api/v1/blog",
			RawQuery: "status=draft",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		reqBody:    `{"title": "Hello, microCMS!"}`,
		resBody:    `{"id": "foo"}`,
	})

	createResult, err := client.Create(CreateParams{
		Endpoint: "blog",
		Status:   StatusDraft,
		Content: Blog{
			Title: "Hello, microCMS!",
		},
	})
	require.NoError(t, err)
	require.EqualValues(
		t,
		&CreateResponse{ID: "foo"},
		createResult,
	)
}

func TestClient_Update(t *testing.T) {
	type Blog struct {
		Title string `json:"title"`
	}

	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodPatch,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog/foo",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
		resHeaders: map[string]string{"Content-Type": "application/json"},
		reqBody:    `{"title": "Hello, microCMS!"}`,
		resBody:    `{"id": "foo"}`,
	})

	updateResult, err := client.Update(UpdateParams{
		Endpoint:  "blog",
		ContentID: "foo",
		Content: Blog{
			Title: "Hello, microCMS!",
		},
	})
	require.NoError(t, err)
	require.EqualValues(
		t,
		&UpdateResponse{ID: "foo"},
		updateResult,
	)
}

func TestClient_Delete(t *testing.T) {
	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(&httpClientMock{
		t:      t,
		method: http.MethodDelete,
		url: &url.URL{
			Scheme: "https",
			Host:   "serviceDomain.microcms.io",
			Path:   "/api/v1/blog/foo",
		},
		reqHeaders: map[string]string{"X-MICROCMS-API-KEY": "apiKey"},
	})

	require.NoError(t, client.Delete(DeleteParams{
		Endpoint:  "blog",
		ContentID: "foo",
	}))
}
