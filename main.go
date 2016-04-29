package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vtex/apps-registrator/models"
	"github.com/vtex/go-gallery-sdk/app"
	"github.com/vtex/go-sdk/vtexid"
	"github.com/vtex/go-sdk/vtexid/apptoken"
)

var (
	appToken   string
	authToken  string
	pathRegexp *regexp.Regexp
)

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

	pathRegexp = regexp.MustCompile("routes/(.*)\\.json")
}

func main() {
	fmt.Println("Listening on port 8080...")
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

	for _, addition := range hook.Additions {
		ConfigureRoutes(hook.Account, hook.Workspace, addition)
	}

	response := models.PreSetupResponse{
		Continue: true,
	}
	responseBody, _ := json.Marshal(response)

	fmt.Fprintf(ctx, string(responseBody))
}

func ConfigureRoutes(account, workspace, addition string) {

	split := strings.Split(addition, "@")
	appId := split[0]
	version := split[1]
	split = strings.Split(appId, ".")
	owner := split[0]
	name := split[1]

	appClient := app.NewAppClient("http://apps.vtex.com", authToken)
	appClient.SetOwner(owner)
	appClient.SetName(name)
	appClient.SetAccount(account)
	appClient.SetWorkspace(workspace)
	appClient.SetAppVersion(version)

	servicesChan, errChan := appClient.ListWorkspaceAppFiles()
	err := <-errChan
	if err != nil {
		panic(err)
	}
	services := <-servicesChan

	for _, service := range services {
		if service == "router" {
			filesChan, errChan := appClient.ListServiceFiles(service)
			err = <-errChan
			if err != nil {
				panic(err)
			}
			files := <-filesChan
			for _, file := range files.Data {
				isCompliant := pathRegexp.MatchString(file.Path)
				if isCompliant {
					bodyChan, errChan := appClient.GetServiceFile(service, file.Path)
					err = <-errChan
					if err != nil {
						panic(err)
					}
					body := <-bodyChan
					route := models.Route{}
					err := json.Unmarshal(body, &route)
					if err != nil {
						panic(err)
					}
					fmt.Println(route.Path)
					fmt.Println(route.Backend)
				}
			}
		}
	}

}
