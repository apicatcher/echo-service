package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 测试 HTTP Echo 接口的功能
func TestHTTPEcho(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := &echoRestAPI{}
	engine.Any("/echo", api.echo)
	t.Run("JSONBody", func(t *testing.T) {
		body := `{"message":"hello","count":123}`
		req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		engine.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data, exists := result["data"]
		if !exists {
			t.Fatal("response does not contain data field")
		}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			t.Fatalf("expected data field to be map, got %T", data)
		}
		if dataMap["message"] != "hello" || dataMap["count"] != float64(123) {
			t.Errorf("unexpected data returned: %v", dataMap)
		}
	})
	t.Run("QueryParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/echo?name=test&val=1", nil)
		resp := httptest.NewRecorder()
		engine.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		data, exists := result["data"]
		if !exists {
			t.Fatal("response does not contain data field")
		}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			t.Fatalf("expected data field to be map, got %T", data)
		}
		if dataMap["name"] != "test" || dataMap["val"] != "1" {
			t.Errorf("unexpected query parameters returned: %v", dataMap)
		}
	})
}
