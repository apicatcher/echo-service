package restapi

import (
	"encoding/json"
	"io"

	"github.com/apicatcher/echo-service/internal/web"
	"github.com/apicatcher/echo-service/internal/web/render"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type echoRestAPI struct {
}

// 自动注册 HTTP Echo 接口路由
func init() {
	echoAPI := &echoRestAPI{}
	web.SetOptions(func(engine *gin.Engine) {
		engine.Any("/echo", echoAPI.echo)
	})
}

// 处理 HTTP Echo 请求，将客户端发送的数据写回在 data 字段中
func (api *echoRestAPI) echo(c *gin.Context) {
	var responseData interface{}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf("failed to read request body: %v", err)
		render.Response(c, nil)
		return
	}
	if len(bodyBytes) > 0 {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonMap); err == nil {
			responseData = jsonMap
		} else {
			responseData = string(bodyBytes)
		}
	} else {
		queryMap := make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				queryMap[k] = v[0]
			}
		}
		responseData = queryMap
	}
	render.Response(c, responseData)
}
