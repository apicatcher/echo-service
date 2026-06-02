package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/apicatcher/echo-service/internal/web"
	_ "github.com/apicatcher/echo-service/internal/web/restapi"
)

func main() {
	runtime.GOMAXPROCS(4)
	go http.ListenAndServe(":8384", nil) // pprof
	web.Start()
}
