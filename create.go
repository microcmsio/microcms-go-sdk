package microcms

import (
	"net/http"
	"net/url"
	"path"
)

type CreateParams struct {
	Endpoint  string
	ContentID string
	Status    string
	Content   interface{}
}

type UpdateParams struct {
	Endpoint  string
	ContentID string
	Content   interface{}
}

type DeleteParams struct {
	Endpoint  string
	ContentID string
}

type CreateResponse struct {
	ID string
}

type UpdateResponse struct {
	ID string
}

func (c *Client) Create(p CreateParams) (*CreateResponse, error) {
	var reqMethod, reqPath string
	if len(p.ContentID) > 0 {
		reqMethod = http.MethodPut
		reqPath = path.Join(p.Endpoint, p.ContentID)
	} else {
		reqMethod = http.MethodPost
		reqPath = p.Endpoint
	}

	req, err := makeRequest(c, reqMethod, reqPath, makeCreateQuery(p), p.Content)
	if err != nil {
		return nil, err
	}

	res := new(CreateResponse)
	if err := sendRequest(c, req, res); err != nil {
		return nil, err
	}

	return res, err
}

func (c *Client) Update(p UpdateParams) (*UpdateResponse, error) {
	req, err := makeRequest(c, http.MethodPatch, path.Join(p.Endpoint, p.ContentID), url.Values{}, p.Content)
	if err != nil {
		return nil, err
	}

	res := new(UpdateResponse)
	if err := sendRequest(c, req, res); err != nil {
		return nil, err
	}

	return res, err
}

func (c *Client) Delete(p DeleteParams) error {
	req, err := makeRequest(c, http.MethodDelete, path.Join(p.Endpoint, p.ContentID), url.Values{}, nil)
	if err != nil {
		return err
	}

	if err := sendRequest(c, req, nil); err != nil {
		return err
	}

	return err
}

func makeCreateQuery(p CreateParams) url.Values {
	urlValues := url.Values{}

	if len(p.Status) > 0 {
		urlValues.Set("status", p.Status)
	}

	return urlValues
}
