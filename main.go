package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/vtex/apps-registrator/models"
	"github.com/vtex/go-sdk/vtexid"
	"github.com/vtex/go-sdk/vtexid/apptoken"
)

var appToken string
var authToken string

func init() {
	var err error
	appToken, err = apptoken.GetValidAppToken()
	if err != nil {
		panic(err)
	}
	authToken, err = vtexid.GetAuthToken("vtexappkey-appvtex", appToken)
	if err != nil {
		panic(err)
	}
}

func main() {
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	hook := models.SetupHook{}
	err := json.Unmarshal(ctx.PostBody(), &hook)
	if err != nil {
		log.Printf("[%s]: %s %s: %s", ctx.RemoteAddr(), method, path, err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		fmt.Fprintf(ctx, "%s", err)
		return
	}

	log.Printf("[%s]: %s %s", ctx.RemoteAddr(), method, path)

	response := models.PreSetupResponse{
		Continue: true,
	}
	responseBody, _ := json.Marshal(response)

	fmt.Fprintf(ctx, string(responseBody))
}
