// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFillStruct_QueryParameters(t *testing.T) {
	e := echo.New()

	t.Run("[]string comma-separated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=value1,value2,value3", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param []string `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value1", "value2", "value3"}, params.Param)
	})

	t.Run("[]string repeated params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=value1&param1=value2&param1=value3", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param []string `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value1", "value2", "value3"}, params.Param)
	})

	t.Run("string", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=value1", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param string `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, "value1", params.Param)
	})

	t.Run("int", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=42", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param int `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, 42, params.Param)
	})

	t.Run("bool", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=true", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param bool `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, true, params.Param)
	})

	t.Run("float64", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?param1=42.42", nil)
		ctx := e.NewContext(req, nil)

		type Request struct {
			Param float64 `in:"query=param1"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, 42.42, params.Param)
	})
}

func TestFillStruct_PathParameters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/path/value1/value2", nil)
	ctx := echo.New().NewContext(req, nil)
	ctx.SetParamNames("param1", "param2")
	ctx.SetParamValues("value1", "value2")

	type Request struct {
		Param1 string `in:"path=param1"`
		Param2 string `in:"path=param2"`
	}

	params := &Request{}
	err := FillStruct(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, "value1", params.Param1)
	assert.Equal(t, "value2", params.Param2)
}

func TestFillStruct_Body(t *testing.T) {
	e := echo.New()

	t.Run("valid JSON", func(t *testing.T) {
		bodyJSON := `{"name":"John Doe"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, nil)

		type Request struct {
			Name string `in:"body=name"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", params.Name)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		bodyJSON := `{"name":"John Doe"`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, nil)

		type Request struct {
			Name string `in:"body=name"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.Error(t, err)
	})
}

func TestFillStruct_Header(t *testing.T) {
	t.Run("valid header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := echo.New().NewContext(req, nil)

		type Request struct {
			ContentType string `in:"header=Content-Type"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, "application/json", params.ContentType)
	})

	t.Run("missing header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := echo.New().NewContext(req, nil)

		type Request struct {
			ContentType string `in:"header=Content-Type"`
		}

		params := &Request{}
		err := FillStruct(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, "", params.ContentType)
	})
}
