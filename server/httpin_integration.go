// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package server

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin/core"
	"github.com/labstack/echo/v4"
)

// EchoMuxVarsFunc is mux.Vars
type EchoMuxVarsFunc func(*http.Request) map[string]string

// UseEchoRouter registers a new directive executor which can extract values
// from `mux.Vars`, i.e. path variables.
// https://ggicci.github.io/httpin/integrations/echo
//
// Usage:
//
//	import httpin_integration "github.com/ggicci/httpin/integration"
//
//	func init() {
//	    e := echo.New()
//	    httpin_integration.UseEchoRouter("path", e)
//	}
func UseEchoRouter(name string, e *echo.Echo) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&echoMuxVarsExtractor{e}).Execute),
		true,
	)
}

func UseEchoPathRouter(e *echo.Echo) {
	UseEchoRouter("path", e)
}

// echoMuxVarsExtractor is an extractor for mux.Vars
type echoMuxVarsExtractor struct {
	e *echo.Echo
}

func (mux *echoMuxVarsExtractor) Execute(rtm *core.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	c := mux.e.NewContext(req, nil)
	c.SetRequest(req)

	mux.e.Router().Find(req.Method, req.URL.Path, c)

	for _, key := range c.ParamNames() {
		kvs[key] = []string{c.Param(key)}
	}

	extractor := &core.FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
