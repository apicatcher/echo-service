package web

import (
	"strconv"

	"github.com/apicatcher/echo-service/internal/config"
	"github.com/apicatcher/echo-service/internal/web/filter"
	"github.com/apicatcher/echo-service/pkg/util"
	"github.com/gin-gonic/gin"
)

var opts []func(*gin.Engine)

func SetOptions(opt func(engine *gin.Engine)) {
	opts = append(opts, opt)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func Start() {
	engine := gin.New()
	middleware := []gin.HandlerFunc{
		filter.AccessLogFilter(),
	}
	engine.Use(middleware...)
	for _, opt := range opts {
		opt(engine)
	}
	if err := engine.Run(":" + strconv.Itoa(config.C.Server.Port)); err != nil {
		util.AbnormalExit(err)
	}
}
