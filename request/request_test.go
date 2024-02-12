// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestParseRequestWithQueryParametersTyped(t *testing.T) {
	t.Run("[]string", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=value1,value2,value3",
			},
		}

		res := []string{}
		// detect if return type is a slice
		ctx := echo.New().NewContext(req, nil)
		err := GetQueryParam[[]string](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value1", "value2", "value3"}, res)
	})
	t.Run("atlernate []string", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=value1&param1=value2&param1=value3",
			},
		}

		res := []string{}
		// detect if return type is a slice
		ctx := echo.New().NewContext(req, nil)
		err := GetQueryParam[[]string](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value1", "value2", "value3"}, res)
	})
	t.Run("string", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=value1",
			},
		}
		ctx := echo.New().NewContext(req, nil)
		res := ""
		err := GetQueryParam[string](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, "value1", res)
	})
	t.Run("int", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=42",
			},
		}
		ctx := echo.New().NewContext(req, nil)
		res := 0
		err := GetQueryParam[int](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, 42, res)
	})
	t.Run("bool", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=true",
			},
		}
		ctx := echo.New().NewContext(req, nil)
		res := false
		err := GetQueryParam[bool](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, true, res)
	})
	t.Run("float64", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				Path:     "/path/to/resource",
				RawQuery: "param1=42.42",
			},
		}
		ctx := echo.New().NewContext(req, nil)
		res := 0.0
		err := GetQueryParam[float64](ctx, "param1", &res)
		assert.NoError(t, err)
		assert.Equal(t, 42.42, res)
	})
}

func TestParseRequestWithPathParameters(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			Path: "/path/to/resource/:param1/:param2",
		},
	}

	ctx := echo.New().NewContext(req, nil)
	ctx.SetParamNames("param1", "param2")
	ctx.SetParamValues("value1", "value2")

	res := ""
	err := GetPathParams[string](ctx, "param1", &res)
	assert.NoError(t, err)
	assert.Equal(t, "value1", res)
	err = GetPathParams[string](ctx, "param2", &res)
	assert.NoError(t, err)
	assert.Equal(t, "value2", res)
}

func TestGetBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		want       interface{}
		wantError  bool
		statusCode int
	}{
		{
			name:       "Get body with valid JSON",
			body:       `{"name": "John Doe"}`,
			want:       &struct{ Name string }{Name: "John Doe"},
			wantError:  false,
			statusCode: http.StatusOK,
		},
		{
			name:       "Get body with invalid JSON",
			body:       `{"name": "John Doe"`,
			want:       nil,
			wantError:  true,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new request with the test body.
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(test.body)))

			// Create a new echo context.
			e := echo.New()

			body := &struct{ Name string }{}
			// Call the GetBody function.
			err := GetBody(e.NewContext(req, nil), body)

			// Check the results.
			if test.wantError {
				if err == nil {
					t.Errorf("GetBody() did not return an error when it should have")
				}
			} else {
				if err != nil {
					t.Errorf("GetBody() returned an error when it should not have: %v", err)
				}
				if !reflect.DeepEqual(body, test.want) {
					t.Errorf("GetBody() got = %v, want = %v", body, test.want)
				}
			}
		})
	}
}
