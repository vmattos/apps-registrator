package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type SetupHook struct {
	Account   string   `json:"account"`
	Workspace string   `json:"workspace"`
	Hash      string   `json:"hash"`
	Removals  []string `json:"removals"`
	Additions []string `json:"additions"`
}

type PreSetupResponse struct {
	Continue bool `json:"continue"`
}

func main() {
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	hook := SetupHook{}
	err := json.Unmarshal(ctx.PostBody(), &hook)
	if err != nil {
		panic(err)
	}

	log.Printf("[%s]: %s %s", ctx.RemoteAddr(), method, path)

	response := PreSetupResponse{
		Continue: true,
	}
	responseBody, _ := json.Marshal(response)

	fmt.Fprintf(ctx, string(responseBody))
}
