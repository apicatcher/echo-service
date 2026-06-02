package restapi

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 测试 WebSocket Echo 和 Ping-Pong 功能
func TestWebSocketEcho(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := &websocketRestAPI{}
	engine.GET("/websocket", api.handleWebSocket)
	server := httptest.NewServer(engine)
	defer server.Close()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/websocket"
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if resp != nil {
			t.Fatalf("failed to dial: %v, status: %d", err, resp.StatusCode)
		}
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	t.Run("TextEcho", func(t *testing.T) {
		err := conn.WriteMessage(websocket.TextMessage, []byte("hello world"))
		if err != nil {
			t.Fatalf("failed to write: %v", err)
		}
		mType, pld, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if mType != websocket.TextMessage {
			t.Errorf("expected TextMessage, got %d", mType)
		}
		if string(pld) != "hello world" {
			t.Errorf("expected 'hello world', got '%s'", string(pld))
		}
	})
	t.Run("TextPingPong", func(t *testing.T) {
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			t.Fatalf("failed to write ping: %v", err)
		}
		mType, pld, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read pong: %v", err)
		}
		if mType != websocket.TextMessage {
			t.Errorf("expected TextMessage, got %d", mType)
		}
		if string(pld) != "pong" {
			t.Errorf("expected 'pong', got '%s'", string(pld))
		}
	})
	t.Run("BinaryPingPong", func(t *testing.T) {
		err := conn.WriteMessage(websocket.BinaryMessage, []byte("ping"))
		if err != nil {
			t.Fatalf("failed to write binary ping: %v", err)
		}
		mType, pld, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read binary pong: %v", err)
		}
		if mType != websocket.BinaryMessage {
			t.Errorf("expected BinaryMessage, got %d", mType)
		}
		if !bytes.Equal(pld, []byte("pong")) {
			t.Errorf("expected binary 'pong', got '%s'", string(pld))
		}
	})
}
