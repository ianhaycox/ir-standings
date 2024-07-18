package main

//type TopSplit map[int]SplitResult

//var points = []int{25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

/*
func TestTable(t *testing.T) {
	b := files.ReadFile(t, "./2024-2-285-results.json")

	res := make([]results.Result, 0)
	err := json.Unmarshal(b, &res)
	assert.NoError(t, err)

	topSplit := make(TopSplit)

	for i := range res {
		subsessionID := res[i].SessionSplits[0].SubsessionID

		if res[i].Track.TrackID == 18 { // Exclude Road America 500 68075008
			continue
		}

		if res[i].SubsessionID == subsessionID {
			splitResult := SplitResult{
				Track:      res[i].Track,
				CarClasses: res[i].CarClasses,
			}

			for j := range res[i].SessionResults {
				if res[i].SessionResults[j].SimsessionName == "RACE" {
					splitResult.Results = res[i].SessionResults[j].Results
				}
			}

			topSplit[subsessionID] = splitResult
		}
	}

	season := Season{
		PointsByCarClass: make(map[CarClassID]map[CustID]SeasonPoints),
	}

	seasonPointsByClass := make(map[int]map[int]SeasonPoints) // map[CarClassID]map[CustID]

	for subSessionID, v := range topSplit {

		resultsByClass := make(map[int][]results.Results)

		for i := range v.Results {
			resultsByClass[v.Results[i].CarClassID] = append(resultsByClass[v.Results[i].CarClassID], v.Results[i])
		}

		for carClassID := range resultsByClass {
			sorted := make([]results.Results, 0)
			sorted = append(sorted, resultsByClass[carClassID]...)

			sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].FinishPositionInClass < sorted[j].FinishPositionInClass })

			for j := range resultsByClass[carClassID] {
				fmt.Printf("    %-30s %2d %s  %d\n", resultsByClass[carClassID][j].DisplayName, resultsByClass[carClassID][j].FinishPositionInClass, resultsByClass[carClassID][j].CarName, resultsByClass[carClassID][j].LapsComplete)
			}

			winnerLapsComplete := sorted[0].LapsComplete

			for j := range resultsByClass[carClassID] {
				if _, ok := seasonPointsByClass[carClassID]; !ok {
					seasonPointsByClass[carClassID] = make(map[int]SeasonPoints)
				}

				if resultsByClass[carClassID][j].LapsComplete*4 <= winnerLapsComplete*3 {
					fmt.Printf("ignoring result %d %d %d\n", winnerLapsComplete, resultsByClass[carClassID][j].LapsComplete, resultsByClass[carClassID][j].CustID)
					continue
				}

				if _, ok := seasonPointsByClass[carClassID][resultsByClass[carClassID][j].CustID]; !ok {
					seasonPointsByClass[carClassID][resultsByClass[carClassID][j].CustID] = SeasonPoints{
						CustID:      resultsByClass[carClassID][j].CustID,
						DisplayName: resultsByClass[carClassID][j].DisplayName,
					}
				}

				sp := seasonPointsByClass[carClassID][resultsByClass[carClassID][j].CustID]

				sp.FinishingPositions = append(sp.FinishingPositions, resultsByClass[carClassID][j].FinishPositionInClass)
				if resultsByClass[carClassID][j].FinishPositionInClass < len(points) {
					sp.ChampionshipsPoints = append(sp.ChampionshipsPoints, points[resultsByClass[carClassID][j].FinishPositionInClass])
				}

				seasonPointsByClass[carClassID][resultsByClass[carClassID][j].CustID] = sp

				fmt.Printf("    %-30s %2d %s\n", resultsByClass[carClassID][j].DisplayName, resultsByClass[carClassID][j].FinishPositionInClass, resultsByClass[carClassID][j].CarName)
			}
		}
	}

	fmt.Println(len(topSplit))

	bb, err := json.MarshalIndent(seasonPointsByClass, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(bb))

	const bestScores = 9

	type Total struct {
		CustID      int
		DisplayName string
		BestOf      int
	}

	for carClassID := range seasonPointsByClass {
		fmt.Printf("%d\n\n", carClassID)
		championshipPositions := make([]Total, 0)

		for custID := range seasonPointsByClass[carClassID] {
			if custID == 544353 {
				fmt.Println(seasonPointsByClass[carClassID][custID])
			}

			total := Total{
				CustID:      custID,
				DisplayName: seasonPointsByClass[carClassID][custID].DisplayName,
			}

			bestOf := make([]int, 0)

			bestOf = append(bestOf, seasonPointsByClass[carClassID][custID].ChampionshipsPoints...)

			sort.SliceStable(bestOf, func(i, j int) bool { return bestOf[i] > bestOf[j] })

			for i := range bestOf {
				if i < bestScores {
					total.BestOf += bestOf[i]
				}
			}

			championshipPositions = append(championshipPositions, total)
		}

		sort.SliceStable(championshipPositions, func(i, j int) bool { return championshipPositions[i].BestOf > championshipPositions[j].BestOf })

		for i := range championshipPositions {
			fmt.Printf("%d  %d  %s %d\n", i+1, championshipPositions[i].BestOf, championshipPositions[i].DisplayName, championshipPositions[i].CustID)
		}
	}
}

func TestStandings(t *testing.T) {
	t.Run("For a list of series results filter and get the full results for broadcast races", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		ssResults := []searchseries.SearchSeriesResult{}
		results := []results.Result{}

		ir := iracing.NewMockIracingService(ctrl)
		ir.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return(ssResults, nil)
		ir.EXPECT().SeasonBroadcastResults(ctx, ssResults).Return(results, nil)

		_, err := standings(ctx, ir, 2024, 2)
		assert.NoError(t, err)
	})

	t.Run("args should return season info", func(t *testing.T) {
		os.Args = []string{"", "1", "2"}
		y, q, err := args()
		assert.NoError(t, err)
		assert.Equal(t, 1, y)
		assert.Equal(t, 2, q)
	})
}

func TestStandingsErrors(t *testing.T) {
	t.Run("Should return error if SearchSeriesResults fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		ssResults := []searchseries.SearchSeriesResult{}

		ir := iracing.NewMockIracingService(ctrl)
		ir.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return(ssResults, fmt.Errorf("failed search"))

		_, err := standings(ctx, ir, 2024, 2)
		assert.ErrorContains(t, err, "failed search")
	})

	t.Run("hould return error if SeasonBroadcastResults fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		ssResults := []searchseries.SearchSeriesResult{}
		results := []results.Result{}

		ir := iracing.NewMockIracingService(ctrl)
		ir.EXPECT().SearchSeriesResults(ctx, 2024, 2, 285).Return(ssResults, nil)
		ir.EXPECT().SeasonBroadcastResults(ctx, ssResults).Return(results, fmt.Errorf("failed broadcast"))

		_, err := standings(ctx, ir, 2024, 2)
		assert.ErrorContains(t, err, "failed broadcast")
	})

	t.Run("args should return error for non numeric year", func(t *testing.T) {
		os.Args = []string{"", "a", "2"}
		_, _, err := args()
		assert.ErrorContains(t, err, "year")
	})

	t.Run("args should return error for no numeric quarter", func(t *testing.T) {
		os.Args = []string{"", "1", "b"}
		_, _, err := args()
		assert.ErrorContains(t, err, "quarter")
	})

	t.Run("args should return error for insufficient args", func(t *testing.T) {
		os.Args = []string{"", "1"}
		_, _, err := args()
		assert.ErrorContains(t, err, "insufficient")
	})
}
*/
