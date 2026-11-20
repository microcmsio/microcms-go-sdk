package microcms

import (
	"net/http"
	"net/url"
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
	reqMethod := http.MethodPost
	if len(p.ContentID) > 0 {
		reqMethod = http.MethodPut
	}

	req, err := makeRequest(c, reqMethod, p.Endpoint, p.ContentID, makeCreateQuery(p), p.Content)
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
	req, err := makeRequest(c, http.MethodPatch, p.Endpoint, p.ContentID, url.Values{}, p.Content)
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
	req, err := makeRequest(c, http.MethodDelete, p.Endpoint, p.ContentID, url.Values{}, nil)
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
