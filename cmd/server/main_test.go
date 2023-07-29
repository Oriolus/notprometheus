package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T,
	ts *httptest.Server,
	method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlers(t *testing.T) {
	testServer := httptest.NewServer(ChiRouter(defaultConfig))
	defer testServer.Close()

	var testTable = []struct {
		url    string
		method string
		want   string
		status int
	}{
		{"", http.MethodGet, "", http.StatusOK},
		{"/update/counter/cnt/1", http.MethodPost, "", http.StatusOK},
		{"/update/counter/testCounter/100", http.MethodPost, "", http.StatusOK},
		{"/update/counter/testCounter/none", http.MethodPost, "", http.StatusBadRequest},
		{"/value/counter/cnt", http.MethodGet, "1", http.StatusOK},
		{"/value/counter/cnt11", http.MethodGet, "", http.StatusNotFound},
		{"/update/gauge/g/1.0", http.MethodPost, "", http.StatusOK},
		{"/value/gauge/g", http.MethodGet, "1", http.StatusOK},
		{"/value/gauge/g11", http.MethodGet, "", http.StatusNotFound},
		{"/value/gg/g11", http.MethodGet, "", http.StatusNotFound},
		{"/update/gg", http.MethodPost, "404 page not found\n", http.StatusNotFound},
		{"/update/gauge/g", http.MethodPost, "404 page not found\n", http.StatusNotFound},
	}

	for _, v := range testTable {
		resp, get := testRequest(t, testServer, v.method, v.url)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
		resp.Body.Close()
	}
}

func testJSONRequest(t *testing.T,
	ts *httptest.Server,
	method, path, contentType string,
	body []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Add("content-type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestJSONedHandler(t *testing.T) {
	testServer := httptest.NewServer(ChiRouter(defaultConfig))
	defer testServer.Close()

	var testTable = []struct {
		url         string
		method      string
		contentType string
		req         []byte
		want        string
		status      int
	}{
		{"/update", http.MethodPost, "multipart/form-data", []byte("{}"), "", http.StatusBadRequest},
		{"/update", http.MethodPost, "application/json", []byte{1}, "", http.StatusBadRequest},
		{"/update", http.MethodPost, "application/json", []byte(`{"id":"124","type":"counter"}`), "", http.StatusBadRequest},
		{"/update", http.MethodPost, "application/json", []byte(`{"id":"counter_json_1","type":"counter","delta":21}`), `{"id":"counter_json_1","type":"counter","value":21}`, http.StatusOK},
		{"/update", http.MethodPost, "application/json", []byte(`{"id":"gauge_json_1","type":"gauge","value":2.22}`), `{"id":"gauge_json_1","type":"gauge","value":2.22}`, http.StatusOK},
		{"/value", http.MethodGet, "multipart/form-data", []byte(`{}`), "", http.StatusBadRequest},
		{"/value", http.MethodGet, "application/json", []byte{1}, "", http.StatusBadRequest},
		{"/value", http.MethodGet, "application/json", []byte(`{"id":1,"type":"cnt"}`), "", http.StatusBadRequest},
		{"/value", http.MethodGet, "application/json", []byte(`{"id":"counter_json_1","type":"counter"}`), `{"id":"counter_json_1","type":"counter","value":21}`, http.StatusOK},
		{"/value", http.MethodGet, "application/json", []byte(`{"id":"gauge_json_1","type":"gauge"}`), `{"id":"gauge_json_1","type":"gauge","value":2.22}`, http.StatusOK},
	}

	for _, v := range testTable {
		resp, get := testJSONRequest(t, testServer, v.method, v.url, v.contentType, v.req)
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
		resp.Body.Close()
	}
}
