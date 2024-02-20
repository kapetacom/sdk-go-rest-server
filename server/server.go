// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type KapetaServer struct {
	*echo.Echo
}

// New creates a new instance of the KapetaServer with default settings
func NewWithDefaults() *KapetaServer {
	e := echo.New()
	e.Add("GET", "/.kapeta/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	// add skipper to skip logging for health check
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/.kapeta/health"
		},
	}))
	// add recover middleware to recover from panics
	e.Use(middleware.Recover())

	return &KapetaServer{e}
}

// New creates a new instance of the KapetaServer, with no default settings
func New() *KapetaServer {
	return &KapetaServer{echo.New()}
}
