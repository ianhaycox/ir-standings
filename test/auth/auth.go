package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
)

// Get authentication token from iRacing and store in a cookie jar
func main() {
	const (
		numArgs = 2
	)

	flag.Parse()

	if len(flag.Args()) != numArgs {
		log.Fatal("insufficient args")
	}

	ctx := context.Background()
	httpClient := http.DefaultClient
	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	httpClient.Jar = cookiejar.NewCookieJar(cookieStore)

	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")

	auth := api.NewAuthenticationService(flag.Arg(0), flag.Arg(1))
	client := api.NewHTTPClient(cfg)

	ir := iracing.NewIracingService(client, nil, auth)

	err := ir.Authenticate(ctx)
	if err != nil {
		log.Fatal("Can not login: ", err)
	}

	fmt.Println("Logged in")
}
