package restapi

import (
	"bytes"
	"net/http"

	"github.com/apicatcher/echo-service/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// websocket 升级配置，允许所有来源的跨域访问
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type websocketRestAPI struct {
}

// 自动注册接口路由
func init() {
	wsAPI := &websocketRestAPI{}
	web.SetOptions(func(engine *gin.Engine) {
		engine.GET("/websocket", wsAPI.handleWebSocket)
	})
}

// 处理 WebSocket 连接并回显消息
func (api *websocketRestAPI) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()
	logrus.Info("websocket connection established")
	for {
		messageType, messagePayload, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorf("failed to read message: %v", err)
			break
		}
		if messageType == websocket.TextMessage && string(messagePayload) == "ping" {
			logrus.Info("received text ping, responding with text pong")
			if err := conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				logrus.Errorf("failed to write pong message: %v", err)
				break
			}
			continue
		}
		if messageType == websocket.BinaryMessage && bytes.Equal(messagePayload, []byte("ping")) {
			logrus.Info("received binary ping, responding with binary pong")
			if err := conn.WriteMessage(websocket.BinaryMessage, []byte("pong")); err != nil {
				logrus.Errorf("failed to write pong message: %v", err)
				break
			}
			continue
		}
		logrus.Infof("echoing message back: type=%d, length=%d", messageType, len(messagePayload))
		if err := conn.WriteMessage(messageType, messagePayload); err != nil {
			logrus.Errorf("failed to echo message: %v", err)
			break
		}
	}
}
