package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
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
		{"/update/counter/cnt/1", http.MethodPost, "", http.StatusOK},
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
