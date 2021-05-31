package microcms

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Base API endpoint
const (
	BASE_DOMAIN = "microcms.io"
	API_VERSION = "v1"
)

// Support HTTP request
const (
	GET = "GET"
)

type HttpResponse *http.Request
type HttpRequest *http.Request

type Client struct {
	serviceDomain  string
	apiKey         string
	globalDraftKey string
}

func createClient(serviceDomain, apiKey, globalDraftKey string) *Client {
	c := &Client{
		serviceDomain:  serviceDomain,
		apiKey:         apiKey,
		globalDraftKey: globalDraftKey,
	}
	return c
}

func createUrl(serviceDomain, endpoint string) string {
	url := fmt.Sprintf("https://%s.%s/api/%s/%s", serviceDomain, BASE_DOMAIN, API_VERSION, endpoint)

	return url
}

func (c *Client) makeRequest(method, endpoint string) ([]byte, error) {
	url := createUrl(c.serviceDomain, endpoint)
	fmt.Println(url)

	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	if c.globalDraftKey != "" {
		req.Header.Set("X-GLOBAL-DRAFT-KEY", c.globalDraftKey)
	}

	res, _ := http.DefaultClient.Do(req)

	resByte, _ := ioutil.ReadAll(res.Body)

	return resByte, nil
}

func (c *Client) Get(endpoint string) (string, error) {
	d, err := c.makeRequest(GET, endpoint)
	return string(d), err
}
