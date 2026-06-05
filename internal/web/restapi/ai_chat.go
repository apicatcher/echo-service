package restapi

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/apicatcher/echo-service/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type aiChatAPI struct {
}

// 自动注册 AI 对话的 SSE 接口路由
func init() {
	chatAPI := &aiChatAPI{}
	web.SetOptions(func(engine *gin.Engine) {
		engine.POST("/v1/chat/completions", chatAPI.chat)
	})
}

// 处理 AI 对话请求，模拟返回 OpenAI 格式的 SSE 流或 JSON
func (api *aiChatAPI) chat(c *gin.Context) {
	var message string
	var isStream bool
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err == nil && len(bodyBytes) > 0 {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonMap); err == nil {
			if stream, ok := jsonMap["stream"].(bool); ok {
				isStream = stream
			}
			if messages, ok := jsonMap["messages"].([]interface{}); ok && len(messages) > 0 {
				lastMsg := messages[len(messages)-1]
				if msgMap, ok := lastMsg.(map[string]interface{}); ok {
					if content, ok := msgMap["content"].(string); ok {
						message = content
					}
				}
			}
		}
	}
	message = strings.ToLower(message)
	logrus.Infof("received ai chat request with message: %s, stream: %v", message, isStream)
	var replyText string
	if strings.Contains(message, "apicatcher") {
		replyText = "ApiCatcher is a focused HTTP/HTTPS capture and debugging tool for iOS. No PC proxy needed, you can capture packets directly on your iPhone or iPad. It supports HTTPS capture, powerful filters, request replay and rewrite, timed replay, and HAR/API export. Furthermore, it features an AI assistant to help you write Regex, Cron expressions, and scripts to effectively improve efficiency."
	} else {
		replyText = "Hello! I am an AI assistant. You can ask me about the ApiCatcher app!"
	}
	if !isStream {
		responseData := map[string]interface{}{
			"id":      "chatcmpl-mock123",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-4-mock",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": replyText,
					},
					"finish_reason": "stop",
				},
			},
		}
		c.JSON(200, responseData)
		return
	}
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	words := strings.Split(replyText, " ")
	index := 0
	c.Stream(func(w io.Writer) bool {
		if index >= len(words) {
			stopData := map[string]interface{}{
				"id":      "apicatcher-mock-123",
				"object":  "chat.completion.chunk",
				"created": time.Now().Unix(),
				"model":   "apicatcher-mock",
				"choices": []map[string]interface{}{
					{
						"index":         0,
						"delta":         map[string]interface{}{},
						"finish_reason": "stop",
					},
				},
			}
			stopBytes, _ := json.Marshal(stopData)
			c.Writer.Write([]byte("data: " + string(stopBytes) + "\n\n"))
			c.Writer.Write([]byte("data: [DONE]\n\n"))
			return false
		}
		word := words[index]
		if index > 0 {
			word = " " + word
		}
		chunkData := map[string]interface{}{
			"id":      "apicatcher-mock-123",
			"object":  "chat.completion.chunk",
			"created": time.Now().Unix(),
			"model":   "apicatcher-mock",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]interface{}{
						"content": word,
					},
					"finish_reason": nil,
				},
			},
		}
		chunkBytes, _ := json.Marshal(chunkData)
		c.Writer.Write([]byte("data: " + string(chunkBytes) + "\n\n"))
		c.Writer.Flush()
		time.Sleep(100 * time.Millisecond)
		index++
		return true
	})
}
