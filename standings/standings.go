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
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
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

	searchSeriesResults, err := ir.SearchSeriesResults(ctx, 2024, 2, iracing.KamelSeriesID)
	if err != nil {
		log.Fatal("Can not get series results:", err)
	}

	allResults := make(map[int]results.Result)

	if searchSeriesResults.Data.Success {
		for i := range searchSeriesResults.Data.ChunkInfo.ChunkFileNames {
			var ssResults []searchseries.SearchSeriesResult

			url := searchSeriesResults.Data.ChunkInfo.BaseDownloadURL + searchSeriesResults.Data.ChunkInfo.ChunkFileNames[i]

			err := data.Get(ctx, url, &ssResults)
			if err != nil {
				log.Fatal("Can not get search series result:"+url, err)
			}

			for j := range ssResults {

				// TODO Just Saturday 17:00 GMT

				if !ssResults[j].IsBroadcast() {
					continue
				}

				link, err := ir.ResultLink(ctx, ssResults[j].SubsessionID)
				if err != nil {
					log.Fatal("Can not get result link for sub session ID:", ssResults[j].SubsessionID, "", err)
				}

				var res results.Result

				err = data.Get(ctx, link.Link, &res)
				if err != nil {
					log.Fatal("Can not get result:"+link.Link, err)
				}

				allResults[ssResults[j].SubsessionID] = res
			}
		}
	}

	b, err := json.MarshalIndent(allResults, "", "  ")
	if err != nil {
		log.Fatal("Can not marshal result:", err.Error())
	}

	err = os.WriteFile("./2024-2-285-results.json", b, 0600) //nolint:mnd // ok
	if err != nil {
		log.Fatal("Can not write result:", err.Error())
	}
}
