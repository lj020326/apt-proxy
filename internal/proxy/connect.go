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
	"fmt"
	"io"
	"net"

	"github.com/gofiber/fiber/v3"
)

// HandleConnect tunnels TLS/TCP via HTTP CONNECT, gated by the same
// passthrough allowlist as HTTPS///. CONNECT is never cached: the proxy only
// sees ciphertext.
func (ap *PackageStruct) HandleConnect(c fiber.Ctx) error {
	targetAddr := string(c.Request().Header.Peek("Host"))
	if targetAddr == "" {
		targetAddr = string(c.Request().URI().Host())
	}
	if targetAddr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing target address")
	}

	hostOnly, _, err := net.SplitHostPort(targetAddr)
	if err != nil {
		hostOnly = targetAddr
		targetAddr = net.JoinHostPort(hostOnly, "443")
	}

	if ap == nil || ap.passthrough == nil || ap.passthrough.Empty() {
		return c.Status(fiber.StatusForbidden).SendString(
			"CONNECT tunneling is disabled (no passthrough origins configured)")
	}

	// Match on host (and optional port) the same way HTTPS/// does.
	if _, allowed := ap.passthrough.Match(hostOnly); !allowed {
		if ap.log != nil {
			ap.log.Warn().
				Str("origin", targetAddr).
				Str("client_ip", c.IP()).
				Msg("CONNECT tunneling rejected: origin not in passthrough allowlist")
		}
		return c.Status(fiber.StatusForbidden).SendString(fmt.Sprintf(
			"%s is not in apt-proxy's passthrough allowlist; "+
				"add it (--passthrough=%s) to allow CONNECT tunneling",
			hostOnly, hostOnly))
	}

	// Dial before hijacking so dial failures can still return a proper status.
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).SendString(
			"502 Bad Gateway: Upstream unreachable")
	}

	fctx := c.RequestCtx()
	fctx.HijackSetNoResponse(true)
	fctx.Hijack(func(clientConn net.Conn) {
		defer func() { _ = clientConn.Close() }()
		defer func() { _ = targetConn.Close() }()

		if _, err := clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
			return
		}

		errChan := make(chan error, 2)
		go func() {
			_, copyErr := io.Copy(targetConn, clientConn)
			if tc, ok := targetConn.(*net.TCPConn); ok {
				_ = tc.CloseWrite()
			}
			errChan <- copyErr
		}()
		go func() {
			_, copyErr := io.Copy(clientConn, targetConn)
			if cc, ok := clientConn.(*net.TCPConn); ok {
				_ = cc.CloseWrite()
			}
			errChan <- copyErr
		}()
		<-errChan
	})

	return nil
}
