package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
)

func main() {
	ctx := context.Background()
	httpClient := http.DefaultClient
	cookieStore := cookiejar.NewStore(iracing.CookiesFile)
	httpClient.Jar = cookiejar.NewCookieJar(cookieStore)

	cfg := api.NewConfiguration(httpClient, api.UserAgent)
	cfg.AddDefaultHeader("Accept", "application/json")
	cfg.AddDefaultHeader("Content-Type", "application/json")
	client := api.NewHTTPClient(cfg)

	// See auth.go that authenticates separately and saves encrypted credentials in a cookie jar
	ir := iracing.NewIracingService(
		client,
		iracing.NewIracingDataService(
			client, cdn.NewCDNService(api.NewHTTPClient(api.NewConfiguration(http.DefaultClient, ""))),
		),
		nil,
	)

	seasonResults, err := ir.Seasons(ctx)
	if err != nil {
		log.Fatal("Can not get results:", err.Error())
	}

	b, err := json.MarshalIndent(seasonResults, "", "  ")
	if err != nil {
		log.Fatal("Can not marshal result:", err.Error())
	}

	fmt.Println(string(b))
}
