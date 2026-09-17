package proxy

import (
	"io"
	"net"

	"github.com/gofiber/fiber/v3"
)

// HandleConnect handles HTTP CONNECT requests to tunnel TLS/TCP traffic.
func HandleConnect(c fiber.Ctx) error {
	// 1. Get host:port from Request URI or Host header
	targetAddr := string(c.Request().Header.Peek("Host"))
	if targetAddr == "" {
		targetAddr = string(c.Request().URI().Host())
	}
	if targetAddr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing target address")
	}

	fctx := c.RequestCtx()

	// 2. Prevent fasthttp from writing its own response before the hijack runs.
	fctx.HijackSetNoResponse(true)

	fctx.Hijack(func(clientConn net.Conn) {
		defer func() { _ = clientConn.Close() }()

		// 3. Dial target upstream server
		targetConn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			_, _ = clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			return
		}
		defer func() { _ = targetConn.Close() }()

		// 4. Send 200 Connection Established back to client
		if _, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
			return
		}

		// 5. Bidirectional copy
		errChan := make(chan error, 2)
		go func() {
			_, err := io.Copy(targetConn, clientConn)
			if tc, ok := targetConn.(*net.TCPConn); ok {
				_ = tc.CloseWrite()
			}
			errChan <- err
		}()

		go func() {
			_, err := io.Copy(clientConn, targetConn)
			if cc, ok := clientConn.(*net.TCPConn); ok {
				_ = cc.CloseWrite()
			}
			errChan <- err
		}()

		<-errChan
	})

	return nil
}
