package microcms

import (
	"fmt"
	"net/http"
)

type HttpResponse *http.Request
type HttpRequest *http.Request

type Client struct {
	serviceDomain  string
	apiKey         string
	globalDraftKey string
}

type ClientParams func(*Client)

func CreateClient(serviceDomain, apiKey string, params ...ClientParams) *Client {
	c := &Client{
		serviceDomain:  serviceDomain,
		apiKey:         apiKey,
		globalDraftKey: "",
	}
	for _, param := range params {
		param(c)
	}
	return c
}

func (c *Client) makeRequest(method, endpoint string) (*http.Request, error) {
	url := createUrl(c.serviceDomain, endpoint)

	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	if c.globalDraftKey != "" {
		req.Header.Set("X-GLOBAL-DRAFT-KEY", c.globalDraftKey)
	}

	return req, nil
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	req, err := c.makeRequest(GET, endpoint)
	res, _ := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, err
}

func GlobalDraftKey(v string) ClientParams {
	return func(c *Client) {
		c.globalDraftKey = v
	}
}

func createUrl(serviceDomain, endpoint string) string {
	url := fmt.Sprintf("https://%s.%s/api/%s/%s", serviceDomain, BASE_DOMAIN, API_VERSION, endpoint)

	return url
}
