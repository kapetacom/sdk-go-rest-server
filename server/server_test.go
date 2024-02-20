// Copyright 2023 Kapeta Inc.
// SPDX-License-Identifier: MIT
package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	// Create a new instance of the KapetaServer with default settings
	// and check if it is not nil
	s := NewWithDefaults()
	assert.NotNil(t, s)
	assert.NotNil(t, s.Echo)
}
