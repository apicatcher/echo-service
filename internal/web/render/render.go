package render

import (
	"github.com/apicatcher/echo-service/internal/dto"
	"github.com/gin-gonic/gin"
)

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	response := dto.NewSuccessResponse(data)
	c.JSON(200, response)
}

// Response 兼容旧版本的响应
func Response(c *gin.Context, data interface{}) {
	ResponseSuccess(c, data)
}
