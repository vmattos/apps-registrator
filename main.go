package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	log.Printf("[%s]: %s %s", ctx.RemoteAddr(), method, path)

	ctx.SetStatusCode(http.StatusOK)
}
