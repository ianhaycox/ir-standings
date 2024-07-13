package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
)

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
	httpClient.Jar = cookiejar.NewCookieJar(iracing.CookiesFile)

	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")

	ir := iracing.NewIracingService(api.NewAPIClient(cfg), nil) // auth.go authenticates separately and saves encrypted credentials in a cookie jar

	data := cdn.NewCDNService(api.NewAPIClient(api.NewConfiguration(http.DefaultClient, "")))

	// https://members-ng.iracing.com/racing/results-stats/results?subsessionid=69999199

	var sessions = []string{"69999199", "70062129", "69930471"}

	for _, sessionID := range sessions {
		link, err := ir.GetResultLink(ctx, sessionID)
		if err != nil {
			log.Fatal("Can not get result link for session ID:", sessionID, "", err)
		}

		result, err := data.GetResult(ctx, link.Link)
		if err != nil {
			log.Fatal("Can not get result:"+link.Link, err)
		}

		fmt.Println(result[0:50])
	}
}
