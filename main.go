package main

import (
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

func main() {
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	fmt.Println(method + " " + path)

	ctx.SetStatusCode(http.StatusOK)
}
