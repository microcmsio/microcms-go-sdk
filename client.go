package microcms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	serviceDomain string
	apiKey        string
	httpClient    httpClient
}

func New(serviceDomain, apiKey string) *Client {
	c := &Client{
		serviceDomain: serviceDomain,
		apiKey:        apiKey,
		httpClient:    http.DefaultClient,
	}
	return c
}

func (c *Client) SetHTTPClient(client httpClient) {
	c.httpClient = client
}

func validatePathSegment(name, segment string) error {
	// Also check every percent-decoded form, so a server or proxy that decodes
	// the path (possibly more than once) cannot turn the segment into a climb.
	v := segment
	for {
		if v == "" || v == "." || v == ".." || strings.ContainsAny(v, `/\`) {
			return fmt.Errorf("microcms: invalid %s %q", name, segment)
		}
		unescaped, err := url.PathUnescape(v)
		if err != nil || unescaped == v {
			return nil
		}
		v = unescaped
	}
}

// makeRequest builds a request to /api/{APIVersion}/{endpoint}[/{contentID}].
// Each identifier is sent as exactly one percent-encoded path segment.
// An empty contentID omits that segment.
func makeRequest(c *Client, method, endpoint, contentID string, query url.Values, data interface{}) (*http.Request, error) {
	if err := validatePathSegment("endpoint", endpoint); err != nil {
		return nil, err
	}
	segments := []string{"api", APIVersion, endpoint}
	if contentID != "" {
		if err := validatePathSegment("content ID", contentID); err != nil {
			return nil, err
		}
		segments = append(segments, contentID)
	}

	escaped := make([]string, len(segments))
	for i, s := range segments {
		escaped[i] = url.PathEscape(s)
	}

	host := fmt.Sprintf("%s.%s", c.serviceDomain, BaseDomain)
	u := &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     "/" + strings.Join(segments, "/"),
		RawPath:  "/" + strings.Join(escaped, "/"),
		RawQuery: query.Encode(),
	}

	buf := new(bytes.Buffer)
	if data != nil {
		if err := json.NewEncoder(buf).Encode(data); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if req.URL.Host != host || req.URL.EscapedPath() != u.EscapedPath() || req.URL.RawQuery != u.RawQuery || req.URL.Fragment != "" {
		return nil, fmt.Errorf("microcms: invalid service domain %q", c.serviceDomain)
	}

	req.Header.Set("X-MICROCMS-API-KEY", c.apiKey)

	if data != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
}

func sendRequest(c *Client, req *http.Request, data interface{}) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() //nolint:errcheck

	if res.StatusCode >= 400 {
		errorMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("microCMS connection error: %w", err)
		}
		return &HttpResponseError{
			Response:     res,
			ErrorMessage: string(errorMessage),
		}
	}

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		if err := json.NewDecoder(res.Body).Decode(data); err != nil {
			return err
		}
	}

	return nil
}

type HttpResponseError struct {
	Response     *http.Response
	ErrorMessage string
}

func (r *HttpResponseError) Error() string {
	return fmt.Sprintf("response error: StatusCode=%d %s", r.Response.StatusCode, r.ErrorMessage)
}
