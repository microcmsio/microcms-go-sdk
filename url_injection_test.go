package microcms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// captureClient records the request makeRequest hands to the HTTP client.
// A 200 with no JSON content type lets sendRequest return without decoding.
type captureClient struct {
	req    *http.Request
	called bool
}

func (c *captureClient) Do(req *http.Request) (*http.Response, error) {
	c.called = true
	c.req = req
	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	return rec.Result(), nil
}

func newClient(t *testing.T) (*Client, *captureClient) {
	t.Helper()
	cap := &captureClient{}
	client := New("serviceDomain", "apiKey")
	client.SetHTTPClient(cap)
	return client, cap
}

// requireRejected demands the call fail before any HTTP request is issued.
// Sending the request and then returning an error is still a failure: the
// API key would already have gone out on that request.
func requireRejected(t *testing.T, err error, cap *captureClient) {
	t.Helper()
	if !cap.called && err != nil {
		return
	}
	target := "<no request>"
	keySet := false
	if cap.called && cap.req != nil && cap.req.URL != nil {
		target = cap.req.Method + " " + cap.req.URL.Scheme + "://" + cap.req.URL.Host + cap.req.URL.RequestURI()
		if cap.req.URL.Fragment != "" {
			target += " fragment=" + cap.req.URL.Fragment
		}
		keySet = cap.req.Header.Get("X-MICROCMS-API-KEY") != ""
	}
	t.Fatalf("expected the client to reject the identifier without an HTTP request; err=%v; request=%s; api key header set=%v", err, target, keySet)
}

// requireNoDotEscape allows a request only when its decoded path stays inside
// the selected endpoint and has no "." or ".." segment. Percent-encoded dot
// segments count, because url.Parse decodes them into URL.Path before the
// request is written.
func requireNoDotEscape(t *testing.T, err error, cap *captureClient, endpoint string) {
	t.Helper()
	if err != nil {
		if cap.called {
			t.Fatalf("request was sent before the client rejected the identifier: %s %s", cap.req.Method, cap.req.URL.RequestURI())
		}
		return
	}
	require.True(t, cap.called, "expected a request or a rejection")
	require.Equal(t, "https", cap.req.URL.Scheme)
	require.Equal(t, "serviceDomain.microcms.io", cap.req.URL.Host)
	require.Empty(t, cap.req.URL.Fragment)

	cleaned := path.Clean(cap.req.URL.Path)
	prefix := "/api/v1/" + endpoint
	if cleaned != prefix && !strings.HasPrefix(cleaned, prefix+"/") {
		t.Fatalf("resolved path %q left endpoint %q (wire %s)", cleaned, endpoint, cap.req.URL.RequestURI())
	}
	for _, seg := range strings.Split(cap.req.URL.Path, "/") {
		if seg == "." || seg == ".." {
			t.Fatalf("decoded path %q contains dot segment %q (wire %s)", cap.req.URL.Path, seg, cap.req.URL.RequestURI())
		}
	}
}

func safePath(segments ...string) string {
	var b strings.Builder
	b.WriteString("/api/v1")
	for _, seg := range segments {
		b.WriteByte('/')
		b.WriteString(url.PathEscape(seg))
	}
	return b.String()
}

func encodeQuery(pairs ...string) string {
	v := url.Values{}
	for i := 0; i < len(pairs); i += 2 {
		v.Set(pairs[i], pairs[i+1])
	}
	return v.Encode()
}

func requireSafeRequest(t *testing.T, err error, cap *captureClient, method, escapedPath, rawQuery string) {
	t.Helper()
	require.NoError(t, err)
	require.True(t, cap.called)
	require.Equal(t, method, cap.req.Method)
	require.Equal(t, "https", cap.req.URL.Scheme)
	require.Equal(t, "serviceDomain.microcms.io", cap.req.URL.Host)
	require.Empty(t, cap.req.URL.Fragment)
	require.Equal(t, escapedPath, cap.req.URL.EscapedPath())
	require.Equal(t, rawQuery, cap.req.URL.RawQuery)
	require.Equal(t, "apiKey", cap.req.Header.Get("X-MICROCMS-API-KEY"))
}

func TestPathSegmentsThatClimbOrAreEmptyAreRejected(t *testing.T) {
	content := map[string]string{"title": "example"}
	out := map[string]interface{}{}

	tests := []struct {
		name string
		call func(c *Client) error
	}{
		{
			name: "get content id climbs out of endpoint",
			call: func(c *Client) error {
				return c.Get(GetParams{
					Endpoint:  "blog",
					ContentID: "../private-endpoint/secret-item",
					DraftKey:  "abcd",
					Fields:    []string{"title"},
				}, &out)
			},
		},
		{
			name: "get content id climbs through an intermediate segment",
			call: func(c *Client) error {
				return c.Get(GetParams{
					Endpoint:  "blog",
					ContentID: "foo/../../private-endpoint/secret-item",
					Fields:    []string{"title"},
				}, &out)
			},
		},
		{
			name: "get content id is dotdot",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog", ContentID: ".."}, &out)
			},
		},
		{
			name: "get content id is dot",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog", ContentID: "."}, &out)
			},
		},
		{
			name: "get content id has an empty segment",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog", ContentID: "foo//bar"}, &out)
			},
		},
		{
			name: "get content id starts with a slash",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog", ContentID: "/private-endpoint/secret-item"}, &out)
			},
		},
		{
			name: "get content id has a trailing slash",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog", ContentID: "secret-item/"}, &out)
			},
		},
		{
			name: "get endpoint climbs and content id is safe",
			call: func(c *Client) error {
				return c.Get(GetParams{
					Endpoint:  "blog/../private-endpoint",
					ContentID: "secret-item",
					DraftKey:  "abcd",
				}, &out)
			},
		},
		{
			name: "get endpoint is dotdot",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "..", ContentID: "secret-item"}, &out)
			},
		},
		{
			name: "get endpoint is dot",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: ".", ContentID: "item"}, &out)
			},
		},
		{
			name: "get object endpoint climbs when content id is empty",
			call: func(c *Client) error {
				return c.Get(GetParams{Endpoint: "blog/../private-endpoint"}, &out)
			},
		},
		{
			name: "update content id climbs out of endpoint",
			call: func(c *Client) error {
				_, err := c.Update(UpdateParams{
					Endpoint:  "blog",
					ContentID: "../private-endpoint/secret-item",
					Content:   content,
				})
				return err
			},
		},
		{
			name: "update content id is dotdot",
			call: func(c *Client) error {
				_, err := c.Update(UpdateParams{Endpoint: "blog", ContentID: "..", Content: content})
				return err
			},
		},
		{
			name: "update endpoint climbs",
			call: func(c *Client) error {
				_, err := c.Update(UpdateParams{
					Endpoint:  "blog/../private-endpoint",
					ContentID: "secret-item",
					Content:   content,
				})
				return err
			},
		},
		{
			name: "delete content id climbs out of endpoint",
			call: func(c *Client) error {
				return c.Delete(DeleteParams{
					Endpoint:  "blog",
					ContentID: "../private-endpoint/secret-item",
				})
			},
		},
		{
			name: "delete content id is dot",
			call: func(c *Client) error {
				return c.Delete(DeleteParams{Endpoint: "blog", ContentID: "."})
			},
		},
		{
			name: "delete endpoint climbs",
			call: func(c *Client) error {
				return c.Delete(DeleteParams{Endpoint: "..", ContentID: "secret-item"})
			},
		},
		{
			name: "create with content id climbs out of endpoint",
			call: func(c *Client) error {
				_, err := c.Create(CreateParams{
					Endpoint:  "blog",
					ContentID: "../private-endpoint/secret-item",
					Status:    StatusDraft,
					Content:   content,
				})
				return err
			},
		},
		{
			name: "create with content id is dotdot",
			call: func(c *Client) error {
				_, err := c.Create(CreateParams{
					Endpoint:  "blog",
					ContentID: "..",
					Content:   content,
				})
				return err
			},
		},
		{
			name: "create without content id uses a climbing endpoint",
			call: func(c *Client) error {
				_, err := c.Create(CreateParams{
					Endpoint: "../private-endpoint",
					Status:   StatusDraft,
					Content:  content,
				})
				return err
			},
		},
		{
			name: "create endpoint climbs through the content id",
			call: func(c *Client) error {
				_, err := c.Create(CreateParams{
					Endpoint:  "blog/../../private-endpoint",
					ContentID: "secret-item",
					Content:   content,
				})
				return err
			},
		},
		{
			name: "list endpoint climbs out of api prefix",
			call: func(c *Client) error {
				return c.List(ListParams{
					Endpoint: "../private-endpoint",
					Filters:  "publishedAt[greater_than]2021-01-01",
					Limit:    5,
					Fields:   []string{"title"},
				}, &out)
			},
		},
		{
			name: "list endpoint climbs through the selected name",
			call: func(c *Client) error {
				return c.List(ListParams{
					Endpoint: "blog/../../private-endpoint",
					Limit:    5,
				}, &out)
			},
		},
		{
			name: "list endpoint is dotdot",
			call: func(c *Client) error {
				return c.List(ListParams{Endpoint: "..", Limit: 5}, &out)
			},
		},
		{
			name: "list endpoint is dot",
			call: func(c *Client) error {
				return c.List(ListParams{Endpoint: ".", Limit: 5}, &out)
			},
		},
		{
			name: "list endpoint has an empty segment",
			call: func(c *Client) error {
				return c.List(ListParams{Endpoint: "blog//secret", Limit: 5}, &out)
			},
		},
		{
			name: "list endpoint has a trailing slash",
			call: func(c *Client) error {
				return c.List(ListParams{Endpoint: "blog/", Filters: "id[equals]foo"}, &out)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cap := newClient(t)
			requireRejected(t, tt.call(client), cap)
		})
	}
}

func TestPercentEncodedDotSegmentsDoNotEscapeEndpoint(t *testing.T) {
	content := map[string]string{"title": "example"}
	out := map[string]interface{}{}

	tests := []struct {
		name     string
		endpoint string
		call     func(c *Client) error
	}{
		{
			name:     "get content id is percent-encoded dotdot",
			endpoint: "blog",
			call: func(c *Client) error {
				return c.Get(GetParams{
					Endpoint:  "blog",
					ContentID: "%2e%2e",
					Fields:    []string{"title"},
				}, &out)
			},
		},
		{
			name:     "get content id is percent-encoded climb",
			endpoint: "blog",
			call: func(c *Client) error {
				return c.Get(GetParams{
					Endpoint:  "blog",
					ContentID: "%2e%2e/%2e%2e/private-endpoint/secret-item",
					DraftKey:  "abcd",
				}, &out)
			},
		},
		{
			name:     "update content id is percent-encoded dotdot",
			endpoint: "blog",
			call: func(c *Client) error {
				_, err := c.Update(UpdateParams{
					Endpoint:  "blog",
					ContentID: "%2E%2E",
					Content:   content,
				})
				return err
			},
		},
		{
			name:     "delete content id is percent-encoded dotdot",
			endpoint: "blog",
			call: func(c *Client) error {
				return c.Delete(DeleteParams{Endpoint: "blog", ContentID: "%2e%2e"})
			},
		},
		{
			name:     "create content id is percent-encoded climb",
			endpoint: "blog",
			call: func(c *Client) error {
				_, err := c.Create(CreateParams{
					Endpoint:  "blog",
					ContentID: "%2e%2e/private-endpoint/secret-item",
					Content:   content,
				})
				return err
			},
		},
		{
			name:     "list endpoint is percent-encoded dotdot",
			endpoint: "blog",
			call: func(c *Client) error {
				return c.List(ListParams{Endpoint: "blog/%2e%2e/private-endpoint", Limit: 5}, &out)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, cap := newClient(t)
			requireNoDotEscape(t, tt.call(client), cap, tt.endpoint)
		})
	}
}

func TestIdentifierCannotSwallowQuery(t *testing.T) {
	content := map[string]string{"title": "example"}
	out := map[string]interface{}{}

	t.Run("get content id hash keeps fields draft key and depth", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Get(GetParams{
			Endpoint:  "blog",
			ContentID: "hello#",
			DraftKey:  "abcd",
			Fields:    []string{"id", "title"},
			Depth:     1,
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog", "hello#"),
			encodeQuery("depth", "1", "draftKey", "abcd", "fields", "id,title"),
		)
	})

	t.Run("get content id question mark does not replace query", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Get(GetParams{
			Endpoint:  "blog",
			ContentID: "hello?x=1",
			DraftKey:  "abcd",
			Fields:    []string{"title"},
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog", "hello?x=1"),
			encodeQuery("draftKey", "abcd", "fields", "title"),
		)
	})

	t.Run("get endpoint hash does not drop content id or fields", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Get(GetParams{
			Endpoint:  "blog#",
			ContentID: "foo",
			Fields:    []string{"title"},
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog#", "foo"),
			encodeQuery("fields", "title"),
		)
	})

	t.Run("get endpoint question mark does not replace fields", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Get(GetParams{
			Endpoint:  "blog?x=1",
			ContentID: "foo",
			Fields:    []string{"title"},
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog?x=1", "foo"),
			encodeQuery("fields", "title"),
		)
	})

	t.Run("list endpoint hash keeps filters limit and fields", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.List(ListParams{
			Endpoint: "blog#",
			DraftKey: "abcd",
			Limit:    5,
			Fields:   []string{"id", "title"},
			Filters:  "publishedAt[greater_than]2021-01-01",
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog#"),
			encodeQuery(
				"draftKey", "abcd",
				"fields", "id,title",
				"filters", "publishedAt[greater_than]2021-01-01",
				"limit", "5",
			),
		)
	})

	t.Run("list endpoint question mark does not replace filters", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.List(ListParams{
			Endpoint: "blog?x=1",
			Limit:    5,
			Filters:  "publishedAt[greater_than]2021-01-01",
		}, &out)
		requireSafeRequest(t, err, cap, http.MethodGet,
			safePath("blog?x=1"),
			encodeQuery("filters", "publishedAt[greater_than]2021-01-01", "limit", "5"),
		)
	})

	t.Run("create content id hash keeps draft status", func(t *testing.T) {
		client, cap := newClient(t)
		_, err := client.Create(CreateParams{
			Endpoint:  "blog",
			ContentID: "item#",
			Status:    StatusDraft,
			Content:   content,
		})
		requireSafeRequest(t, err, cap, http.MethodPut,
			safePath("blog", "item#"),
			encodeQuery("status", "draft"),
		)
	})

	t.Run("create content id question mark keeps draft status", func(t *testing.T) {
		client, cap := newClient(t)
		_, err := client.Create(CreateParams{
			Endpoint:  "blog",
			ContentID: "item?x=1",
			Status:    StatusDraft,
			Content:   content,
		})
		requireSafeRequest(t, err, cap, http.MethodPut,
			safePath("blog", "item?x=1"),
			encodeQuery("status", "draft"),
		)
	})

	t.Run("create endpoint hash keeps draft status", func(t *testing.T) {
		client, cap := newClient(t)
		_, err := client.Create(CreateParams{
			Endpoint: "blog#",
			Status:   StatusDraft,
			Content:  content,
		})
		requireSafeRequest(t, err, cap, http.MethodPost,
			safePath("blog#"),
			encodeQuery("status", "draft"),
		)
	})

	t.Run("update content id question mark adds no query", func(t *testing.T) {
		client, cap := newClient(t)
		_, err := client.Update(UpdateParams{
			Endpoint:  "blog",
			ContentID: "item?status=draft",
			Content:   content,
		})
		requireSafeRequest(t, err, cap, http.MethodPatch,
			safePath("blog", "item?status=draft"),
			"",
		)
	})

	t.Run("update content id hash stays in the path", func(t *testing.T) {
		client, cap := newClient(t)
		_, err := client.Update(UpdateParams{
			Endpoint:  "blog",
			ContentID: "item#",
			Content:   content,
		})
		requireSafeRequest(t, err, cap, http.MethodPatch,
			safePath("blog", "item#"),
			"",
		)
	})

	t.Run("delete content id hash stays in the path", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Delete(DeleteParams{Endpoint: "blog", ContentID: "item#"})
		requireSafeRequest(t, err, cap, http.MethodDelete,
			safePath("blog", "item#"),
			"",
		)
	})

	t.Run("delete content id question mark adds no query", func(t *testing.T) {
		client, cap := newClient(t)
		err := client.Delete(DeleteParams{Endpoint: "blog", ContentID: "item?x=1"})
		requireSafeRequest(t, err, cap, http.MethodDelete,
			safePath("blog", "item?x=1"),
			"",
		)
	})
}
