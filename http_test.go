/**
 * Copyright 2024-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRestClient struct {
	baseUrl    string
	httpClient *http.Client
}

func (c *testRestClient) HttpBaseUrl() string      { return c.baseUrl }
func (c *testRestClient) HttpClient() *http.Client { return c.httpClient }

func newTestRestClient(t *testing.T, handler http.HandlerFunc) RestClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &testRestClient{baseUrl: srv.URL, httpClient: srv.Client()}
}

func TestHttpGetSuccess(t *testing.T) {
	cl := newTestRestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	var resp struct {
		OK bool `json:"ok"`
	}
	if err := HttpGet(context.Background(), cl, "/ok", EmptyQueryParams, []int{http.StatusOK}, struct{}{}, &resp, nil, nil); err != nil {
		t.Fatal(err)
	}
	if !resp.OK {
		t.Fatal("expected ok")
	}
}

func TestHttpGetNilParserReturnsApiError(t *testing.T) {
	cl := newTestRestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
	})

	var resp struct{}
	err := HttpGet(context.Background(), cl, "/fail", EmptyQueryParams, []int{http.StatusOK}, struct{}{}, &resp, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T (%v)", err, err)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("message = %q", apiErr.Message)
	}
	if apiErr.CodeReceived != http.StatusBadRequest {
		t.Errorf("status = %d", apiErr.CodeReceived)
	}
}

type customParseError struct {
	msg string
}

func (e *customParseError) Error() string { return e.msg }

func TestHttpGetCustomParserReturnedAsIs(t *testing.T) {
	cl := newTestRestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"NOT_FOUND"}`))
	})

	parser := func(body []byte, statusCode int, expected []int, callUrl string) error {
		return &customParseError{msg: fmt.Sprintf("%d:%s", statusCode, string(body))}
	}

	var resp struct{}
	err := HttpGet(context.Background(), cl, "/x", EmptyQueryParams, []int{http.StatusOK}, struct{}{}, &resp, nil, parser)
	if err == nil {
		t.Fatal("expected error")
	}
	var custom *customParseError
	if !errors.As(err, &custom) {
		t.Fatalf("expected *customParseError, got %T (%v)", err, err)
	}
	if custom.msg != `404:{"code":"NOT_FOUND"}` {
		t.Errorf("msg = %q", custom.msg)
	}
}

func TestAppendHttpQueryParam(t *testing.T) {

	cases := []struct {
		description string
		queryParams string
		key         string
		value       string
		expected    string
	}{
		{
			description: "TestAppendHttpQueryParam0",
			queryParams: "",
			key:         "foo",
			value:       "bar",
			expected:    "?foo=bar",
		},
		{
			description: "TestAppendHttpQueryParam1",
			queryParams: "?test=new",
			key:         "foo",
			value:       "bar",
			expected:    "?test=new&foo=bar",
		},
		{
			description: "TestAppendHttpQueryParam2",
			queryParams: "?test=new&new=test",
			key:         "foo",
			value:       "bar",
			expected:    "?test=new&new=test&foo=bar",
		},
		{
			description: "TestAppendHttpQueryParam3",
			queryParams: "",
			key:         "foo",
			value:       "hello world&bad=injected",
			expected:    "?foo=hello+world%26bad%3Dinjected",
		},
		{
			description: "TestAppendHttpQueryParam4",
			queryParams: "?a=1",
			key:         "foo bar",
			value:       "baz=qux",
			expected:    "?a=1&foo+bar=baz%3Dqux",
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			result := AppendHttpQueryParam(tt.queryParams, tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("test: %s - expected: %s - received: %s", tt.description, tt.expected, result)
			}
		})
	}
}
