package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
	"github.com/ianhaycox/ir-standings/model/championship"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/data/results"
)

func main() {
	seasonYear, seasonQuarter, err := args() // TODO SeriesID, BestOf, Exclude events/tracks
	if err != nil {
		exampleData, err := getResults()
		if err != nil {
			log.Fatal("Can not get example results:", err.Error())
		}

		const maxSplits = 3

		var excludeTrackID = map[int]bool{18: true}

		pointsPerSplit := points.PointsPerSplit{
			0: []int{25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			1: []int{14, 12, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			2: []int{9, 6, 4, 3, 2, 1},
		}

		ps := points.NewPointsStructure(pointsPerSplit)

		champ := championship.NewChampionship(iracing.KamelSeriesID, excludeTrackID, maxSplits, ps)

		champ.LoadRaceData(exampleData)

		return
	}

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

	seasonResults, err := standings(ctx, ir, seasonYear, seasonQuarter)
	if err != nil {
		log.Fatal("Can not get standings:", err.Error())
	}

	b, err := json.MarshalIndent(seasonResults, "", "  ")
	if err != nil {
		log.Fatal("Can not marshal result:", err.Error())
	}

	err = os.WriteFile(fmt.Sprintf("./%d-%d-%d-results.json", seasonYear, seasonQuarter, iracing.KamelSeriesID), b, 0600) //nolint:mnd // ok
	if err != nil {
		log.Fatal("Can not write result:", err.Error())
	}
}

func standings(ctx context.Context, ir iracing.IracingService, seasonYear, seasonQuarter int) ([]results.Result, error) {
	searchSeriesResults, err := ir.SearchSeriesResults(ctx, seasonYear, seasonQuarter, iracing.KamelSeriesID)
	if err != nil {
		return nil, fmt.Errorf("can not get series results:%w", err)
	}

	seasonResults, err := ir.SeasonBroadcastResults(ctx, searchSeriesResults)
	if err != nil {
		return nil, fmt.Errorf("can not get series results:%w", err)
	}

	return seasonResults, nil
}

func args() (int, int, error) {
	const (
		numArgs = 2
	)

	flag.Parse()

	if len(flag.Args()) != numArgs {
		return 0, 0, fmt.Errorf("insufficient args")
	}

	seasonYear, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return 0, 0, fmt.Errorf("season year should be numeric, e.g. 2023")
	}

	seasonQuarter, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		return 0, 0, fmt.Errorf("season quarter should be numeric, e.g. 1, 2, 3, 4, 5")
	}

	return seasonYear, seasonQuarter, nil
}

func getResults() ([]results.Result, error) {
	buf, err := os.ReadFile("./2024-2-285-results.json")
	if err != nil {
		return nil, err
	}

	res := make([]results.Result, 0)

	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
