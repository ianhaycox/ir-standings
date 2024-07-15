package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	result "github.com/ianhaycox/ir-standings/model"
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
	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	httpClient.Jar = cookiejar.NewCookieJar(cookieStore)

	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")

	// See auth.go that authenticates separately and saves encrypted credentials in a cookie jar
	ir := iracing.NewIracingService(iracing.NewIracingDataService(api.NewAPIClient(cfg)), nil)

	data := cdn.NewCDNService(api.NewAPIClient(api.NewConfiguration(http.DefaultClient, "")))

	// https://members-ng.iracing.com/racing/results-stats/results?subsessionid=69999199

	var sessions = []int{69999199, 70062129, 69930471}

	results := make(map[int]result.Result)

	for _, sessionID := range sessions {
		link, err := ir.ResultLink(ctx, sessionID)
		if err != nil {
			log.Fatal("Can not get result link for session ID:", sessionID, "", err)
		}

		var res result.Result

		err = data.Get(ctx, link.Link, &res)
		if err != nil {
			log.Fatal("Can not get result:"+link.Link, err)
		}

		results[sessionID] = res
	}

	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal("Can not marshal result:", err.Error())
	}

	err = os.WriteFile("./r.json", b, 0600) //nolint:mnd // ok
	if err != nil {
		log.Fatal("Can not write result:", err.Error())
	}
}
