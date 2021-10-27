package microcms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	serviceDomain string
	apiKey        string
}

type Params struct {
	contentID string
	draftKey  string
	limit     int
	offset    int
	orders    []string
	q         string
	fields    []string
	ids       []string
	filters   string
	depth     int
}

type RequestParams func(*Params)

func CreateClient(serviceDomain, apiKey string) *Client {
	c := &Client{
		serviceDomain: serviceDomain,
		apiKey:        apiKey,
	}
	return c
}

func (c *Client) makeRequest(method, endpoint string, p *Params) (*http.Request, error) {
	url := createUrl(c.serviceDomain, endpoint, p)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MICROCMS-API-KEY", c.apiKey)

	return req, nil
}

func (c *Client) Get(endpoint string, data interface{}, params ...RequestParams) error {
	p := &Params{}

	for _, params := range params {
		params(p)
	}

	req, err := c.makeRequest(http.MethodGet, endpoint, p)
	res, _ := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := parseBody(res, &data); err != nil {
		return err
	}

	return err
}

func parseBody(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(v)
}

func createUrl(serviceDomain, endpoint string, p *Params) string {
	base := fmt.Sprintf("https://%s.%s/api/%s/%s", serviceDomain, BaseDomain, APIVersion, endpoint)

	if p.contentID != "" {
		base = fmt.Sprintf("%s/%s", base, p.contentID)
	}

	urlValues := url.Values{}
	if len(p.draftKey) > 0 {
		urlValues.Set("draftKey", p.draftKey)
	}
	if p.limit != 0 {
		urlValues.Set("limit", fmt.Sprint(p.limit))
	}
	if p.offset != 0 {
		urlValues.Set("offset", fmt.Sprint(p.offset))
	}
	if len(p.orders) > 0 {
		urlValues.Set("orders", strings.Join(p.orders, ","))
	}
	if len(p.q) > 0 {
		urlValues.Set("q", p.q)
	}
	if len(p.fields) > 0 {
		urlValues.Set("fields", strings.Join(p.fields, ","))
	}
	if len(p.ids) > 0 {
		urlValues.Set("ids", strings.Join(p.ids, ","))
	}
	if len(p.filters) > 0 {
		urlValues.Set("filters", p.filters)
	}
	if p.depth != 0 {
		urlValues.Set("depth", fmt.Sprint(p.depth))
	}

	if len(urlValues) > 0 {
		base = fmt.Sprintf("%s?%s", base, urlValues.Encode())
	}

	return base
}

func ContentID(v string) RequestParams {
	return func(p *Params) {
		p.contentID = v
	}
}

func DraftKey(v string) RequestParams {
	return func(p *Params) {
		p.draftKey = v
	}
}

func Limit(v int) RequestParams {
	return func(p *Params) {
		p.limit = v
	}
}

func Offset(v int) RequestParams {
	return func(p *Params) {
		p.offset = v
	}
}

func Orders(v []string) RequestParams {
	return func(p *Params) {
		p.orders = v
	}
}

func Q(v string) RequestParams {
	return func(p *Params) {
		p.q = v
	}
}

func Fields(v []string) RequestParams {
	return func(p *Params) {
		p.fields = v
	}
}

func IDs(ids []string) RequestParams {
	return func(p *Params) {
		p.ids = ids
	}
}

func Filters(v string) RequestParams {
	return func(p *Params) {
		p.filters = v
	}
}

func Depth(v int) RequestParams {
	return func(p *Params) {
		p.depth = v
	}
}
