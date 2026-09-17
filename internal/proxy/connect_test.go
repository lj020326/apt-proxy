package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleConnect(t *testing.T) {
	tests := []struct {
		name  string
		isTLS bool
	}{
		{
			name:  "HTTP Target Tunnel",
			isTLS: false,
		},
		{
			name:  "HTTPS Target Tunnel",
			isTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello via CONNECT Tunnel")
			})

			var ts *httptest.Server
			if tt.isTLS {
				ts = httptest.NewTLSServer(handler)
			} else {
				ts = httptest.NewServer(handler)
			}
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			// Setup Fiber Proxy Server with CONNECT route
			app := fiber.New()
			app.Add([]string{fiber.MethodConnect}, "*", HandleConnect)

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer listener.Close()

			go func() {
				_ = app.Listener(listener) // or with ListenConfig if your version supports it
			}()

			proxyAddr := listener.Addr().String()

			// Dial proxy using CONNECT method
			proxyConn, err := net.Dial("tcp", proxyAddr)
			require.NoError(t, err)
			defer proxyConn.Close()

			connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", u.Host, u.Host)
			_, err = proxyConn.Write([]byte(connectReq))
			require.NoError(t, err)

			// Read status line without over-buffering bytes from the connection
			bufReader := bufio.NewReader(proxyConn)
			statusLine, err := bufReader.ReadString('\n')
			require.NoError(t, err)
			assert.Contains(t, statusLine, "200")

			// Consume remaining headers until empty line (\r\n)
			for {
				line, err := bufReader.ReadString('\n')
				require.NoError(t, err)
				if strings.TrimSpace(line) == "" {
					break
				}
			}

			// Prepare target connection (plain TCP tunnel or TLS)
			var targetConn = proxyConn
			if tt.isTLS {
				targetConn = tls.Client(proxyConn, &tls.Config{
					InsecureSkipVerify: true,
					ServerName:         u.Hostname(),
				})
				defer targetConn.Close()
			}

			// Send HTTP GET through the tunnel
			req, err := http.NewRequest("GET", ts.URL, nil)
			require.NoError(t, err)
			req.Header.Set("Connection", "close")
			err = req.Write(targetConn)
			require.NoError(t, err)

			targetBuf := bufio.NewReader(targetConn)
			targetRes, err := http.ReadResponse(targetBuf, req)
			require.NoError(t, err)
			defer targetRes.Body.Close()
			assert.Equal(t, http.StatusOK, targetRes.StatusCode)

			body, err := io.ReadAll(targetRes.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), "Hello via CONNECT Tunnel")
		})
	}
}
