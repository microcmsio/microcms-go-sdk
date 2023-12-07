package microcms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func makeRequest(c *Client, method, endpoint string, query url.Values, data interface{}) (*http.Request, error) {
	url := fmt.Sprintf("https://%s.%s/api/%s/%s", c.serviceDomain, BaseDomain, APIVersion, endpoint)
	if len(query) > 0 {
		url = fmt.Sprintf("%s?%s", url, query.Encode())
	}

	buf := new(bytes.Buffer)
	if data != nil {
		if err := json.NewEncoder(buf).Encode(data); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
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
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errorMessage, err := ioutil.ReadAll(res.Body)
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
