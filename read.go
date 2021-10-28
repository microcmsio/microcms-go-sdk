package microcms

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type ListParams struct {
	Endpoint string
	DraftKey string
	Limit    int
	Offset   int
	Orders   []string
	Q        string
	Fields   []string
	IDs      []string
	Filters  string
	Depth    int
}

type GetParams struct {
	Endpoint  string
	ContentID string
	DraftKey  string
	Fields    []string
	Depth     int
}

func (c *Client) List(p ListParams, data interface{}) error {
	req, err := makeRequest(c, http.MethodGet, p.Endpoint, makeListQuery(p), nil)
	if err != nil {
		return err
	}

	if err := sendRequest(c, req, data); err != nil {
		return err
	}

	return err
}

func (c *Client) Get(p GetParams, data interface{}) error {
	req, err := makeRequest(c, http.MethodGet, path.Join(p.Endpoint, p.ContentID), makeGetQuery(p), nil)
	if err != nil {
		return err
	}

	if err := sendRequest(c, req, data); err != nil {
		return err
	}

	return err
}

func makeListQuery(p ListParams) url.Values {
	urlValues := url.Values{}

	if len(p.DraftKey) > 0 {
		urlValues.Set("draftKey", p.DraftKey)
	}
	if p.Limit != 0 {
		urlValues.Set("limit", fmt.Sprint(p.Limit))
	}
	if p.Offset != 0 {
		urlValues.Set("offset", fmt.Sprint(p.Offset))
	}
	if len(p.Orders) > 0 {
		urlValues.Set("orders", strings.Join(p.Orders, ","))
	}
	if len(p.Q) > 0 {
		urlValues.Set("q", p.Q)
	}
	if len(p.Fields) > 0 {
		urlValues.Set("fields", strings.Join(p.Fields, ","))
	}
	if len(p.IDs) > 0 {
		urlValues.Set("ids", strings.Join(p.IDs, ","))
	}
	if len(p.Filters) > 0 {
		urlValues.Set("filters", p.Filters)
	}
	if p.Depth != 0 {
		urlValues.Set("depth", fmt.Sprint(p.Depth))
	}

	return urlValues
}

func makeGetQuery(p GetParams) url.Values {
	urlValues := url.Values{}

	if len(p.DraftKey) > 0 {
		urlValues.Set("draftKey", p.DraftKey)
	}
	if len(p.Fields) > 0 {
		urlValues.Set("fields", strings.Join(p.Fields, ","))
	}
	if p.Depth != 0 {
		urlValues.Set("depth", fmt.Sprint(p.Depth))
	}

	return urlValues
}
