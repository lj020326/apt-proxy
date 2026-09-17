// Copyright 2026 LJ Johnson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectPassthroughRestriction(t *testing.T) {
	// Public hostname only: passthrough.Parse refuses non-public addresses.
	ps, _ := passthroughProxy(t, "example.com")

	app := fiber.New()
	app.Add([]string{fiber.MethodConnect}, "*", ps.HandleConnect)

	tests := []struct {
		name           string
		targetHost     string
		expectedStatus int
	}{
		{
			name:           "unauthorized origin rejected",
			targetHost:     "malicious-packages.example.com:443",
			expectedStatus: http.StatusForbidden,
		},
		{
			// Port 9 is unassigned; dial fails after the allowlist check passes.
			name:           "allowed origin, unreachable port",
			targetHost:     "example.com:9",
			expectedStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodConnect, "//"+tt.targetHost, nil)
			require.NoError(t, err)
			req.Host = tt.targetHost

			resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestConnectDisabledWithoutPassthrough(t *testing.T) {
	ps := markerTestProxy(t, "http://127.0.0.1:9") // no Passthrough set

	app := fiber.New()
	app.Add([]string{fiber.MethodConnect}, "*", ps.HandleConnect)

	req, err := http.NewRequest(http.MethodConnect, "//example.com:443", nil)
	require.NoError(t, err)
	req.Host = "example.com:443"

	resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
