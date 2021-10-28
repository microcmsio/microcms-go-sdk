package microcms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type httpClientMock struct {
	t          *testing.T
	method     string
	url        *url.URL
	reqHeaders map[string]string
	resHeaders map[string]string
	reqBody    string
	resBody    string
}

var _ httpClient = &httpClientMock{}

func (m *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	require.Equal(m.t, m.method, req.Method)
	require.Equal(m.t, m.url, req.URL)

	for k, v := range m.reqHeaders {
		require.Equal(m.t, v, req.Header.Get(k))
	}

	if len(m.reqBody) > 0 {
		var expectedReqBody interface{}
		require.NoError(m.t, json.NewDecoder(bytes.NewBufferString(m.reqBody)).Decode(&expectedReqBody))

		var actualReqBody interface{}
		require.NoError(m.t, json.NewDecoder(req.Body).Decode(&actualReqBody))

		require.EqualValues(m.t, expectedReqBody, actualReqBody)
	}

	rec := httptest.NewRecorder()

	for k, v := range m.resHeaders {
		rec.Header().Set(k, v)
	}
	if _, err := rec.Body.Write([]byte(m.resBody)); err != nil {
		return nil, err
	}

	return rec.Result(), nil
}
