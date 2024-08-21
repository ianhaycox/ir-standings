package predictor

/*

	i83 := []live.PredictedStanding{
		{CustID: 1, DriverName: "Audi 1", CarNumber: "1", CurrentPosition: 1, PredictedPosition: 2, CurrentPoints: 100, PredictedPoints: 100, Change: -1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 2, DriverName: "Audi 2", CarNumber: "2", CurrentPosition: 2, PredictedPosition: 1, CurrentPoints: 95, PredictedPoints: 110, Change: 1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 3, DriverName: "Audi 3", CarNumber: "3", CurrentPosition: 3, PredictedPosition: 3, CurrentPoints: 90, PredictedPoints: 90, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 4, DriverName: "Audi 4", CarNumber: "4", CurrentPosition: 4, PredictedPosition: 4, CurrentPoints: 85, PredictedPoints: 85, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 5, DriverName: "Audi 5", CarNumber: "5", CurrentPosition: 5, PredictedPosition: 5, CurrentPoints: 80, PredictedPoints: 80, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 6, DriverName: "Audi 6", CarNumber: "6", CurrentPosition: 6, PredictedPosition: 7, CurrentPoints: 75, PredictedPoints: 75, Change: -1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 7, DriverName: "Audi 7", CarNumber: "7", CurrentPosition: 7, PredictedPosition: 6, CurrentPoints: 70, PredictedPoints: 78, Change: 1, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 8, DriverName: "Audi 8", CarNumber: "8", CurrentPosition: 8, PredictedPosition: 8, CurrentPoints: 65, PredictedPoints: 65, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 9, DriverName: "Audi 9", CarNumber: "9", CurrentPosition: 9, PredictedPosition: 9, CurrentPoints: 60, PredictedPoints: 60, Change: 0, CarNames: []string{"Audi 90 GTO"}},
		{CustID: 10, DriverName: "Audi 10", CarNumber: "10", CurrentPosition: 10, PredictedPosition: 10, CurrentPoints: 55, PredictedPoints: 55, Change: 0, CarNames: []string{"Audi 90 GTO"}},
	}

	i84 := []live.PredictedStanding{
		{CustID: 11, DriverName: "Nissan 1", CarNumber: "1", CurrentPosition: 1, PredictedPosition: 2, CurrentPoints: 100, PredictedPoints: 100, Change: -1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 12, DriverName: "Nissan 2", CarNumber: "2", CurrentPosition: 2, PredictedPosition: 1, CurrentPoints: 95, PredictedPoints: 110, Change: 1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 13, DriverName: "Nissan 3", CarNumber: "3", CurrentPosition: 3, PredictedPosition: 3, CurrentPoints: 90, PredictedPoints: 90, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 14, DriverName: "Nissan 4", CarNumber: "4", CurrentPosition: 4, PredictedPosition: 4, CurrentPoints: 85, PredictedPoints: 85, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 15, DriverName: "Nissan 5", CarNumber: "5", CurrentPosition: 5, PredictedPosition: 5, CurrentPoints: 80, PredictedPoints: 80, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 16, DriverName: "Nissan 6", CarNumber: "6", CurrentPosition: 6, PredictedPosition: 7, CurrentPoints: 75, PredictedPoints: 75, Change: -1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 17, DriverName: "Nissan 7", CarNumber: "7", CurrentPosition: 7, PredictedPosition: 6, CurrentPoints: 70, PredictedPoints: 78, Change: 1, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 18, DriverName: "Nissan 8", CarNumber: "8", CurrentPosition: 8, PredictedPosition: 8, CurrentPoints: 65, PredictedPoints: 65, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 19, DriverName: "Nissan 9", CarNumber: "9", CurrentPosition: 9, PredictedPosition: 9, CurrentPoints: 60, PredictedPoints: 60, Change: 0, CarNames: []string{"Nissan ZX-T"}},
		{CustID: 20, DriverName: "Nissan 10", CarNumber: "10", CurrentPosition: 10, PredictedPosition: 10, CurrentPoints: 55, PredictedPoints: 55, Change: 0, CarNames: []string{"Nissan ZX-T"}},
	}

	s83 := live.Standing{
		SoFByCarClass:           2087,
		CarClassID:              83,
		CarClassName:            "GTO",
		ClassLeaderLapsComplete: 10,
		Items:                   i83,
	}

	s84 := live.Standing{
		SoFByCarClass:           3025,
		CarClassID:              84,
		CarClassName:            "GTP",
		ClassLeaderLapsComplete: 12,
		Items:                   i84,
	}

	ps := live.PredictedStandings{
		TrackName:   "Motegi Resort",
		CountBestOf: 10,
		Standings:   map[model.CarClassID]live.Standing{84: s84, 83: s83},
	}

func (c *SafeChamp) load(seriesID model.SeriesID, carClassID model.CarClassID, filename string, countBestOf int) error {
	if c.previous == nil {
		c.previous = championship.NewChampionship(seriesID, nil, c.ps, countBestOf)

		buf, err := os.ReadFile(filename) //nolint:gosec // ok
		if err != nil {
			return fmt.Errorf("can not open file %s", filename)
		}

		x := string(buf)
		l := strings.Index(x, "Dustin")

		fmt.Println(x[l-10 : l+20])

		err = json.Unmarshal(buf, &c.previousResults)
		if err != nil {
			return fmt.Errorf("can not parse file %s", filename)
		}

		c.previous.LoadRaceData(c.previousResults)
		c.previousStandings = c.previous.Standings(carClassID)
	}

	return nil
}

func (c *SafeChamp) Live(filename string, jsonCurrentPositions string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ps = points.NewPointsStructure(pointsPerSplit)

	var currentPositions live.LiveResults

	err := json.Unmarshal([]byte(jsonCurrentPositions), &currentPositions)
	if err != nil {
		return "", fmt.Errorf("malformed request %w", err)
	}

	carClassID := model.CarClassID(currentPositions.CarClassID)
	seriesID := model.SeriesID(currentPositions.SeriesID)

	err = c.load(seriesID, carClassID, filename, currentPositions.CountBestOf)
	if err != nil {
		return "", fmt.Errorf("can load previous results, err:%w", err)
	}

	//	f.WriteString(fmt.Sprintf("%s Live:%d, %s\n", time.Now(), carClassID, jsonCurrentPositions))

	liveResults := make([]results.Result, 0, len(c.previousResults))
	liveResults = append(liveResults, c.previousResults...)

	liveResults = append(liveResults, results.Result{
		SessionID:     currentPositions.SessionID,
		SubsessionID:  currentPositions.SubsessionID,
		SeriesID:      currentPositions.SeriesID,
		SessionSplits: []results.SessionSplits{{SubsessionID: currentPositions.SubsessionID}},
		StartTime:     time.Now().UTC(),
		Track:         results.ResultTrack{TrackID: 1, TrackName: "track"},
		SessionResults: []results.SessionResults{
			{
				SimsessionName: "RACE",
				Results:        buildResults(carClassID, currentPositions.Positions),
			},
		},
	})

	predicted := championship.NewChampionship(seriesID, nil, c.ps, currentPositions.CountBestOf)

	predicted.LoadRaceData(liveResults)
	predicted.SetCarClasses(c.previous.CarClasses())

	predictedStandings := predicted.Standings(carClassID)

	//	f.WriteString(fmt.Sprintf("%s Predicted:%d %+v\n\n", time.Now(), carClassID, predictedStandings))

	currentStandings := c.previous.Standings(carClassID)

	provisionalChampionship := provisionalTable(currentStandings, predictedStandings)

	jsonResult, err := json.MarshalIndent(provisionalChampionship[:currentPositions.TopN], "", "  ")
	if err != nil {
		return "", fmt.Errorf("can not build response, %w", err)
	}

	//	f.WriteString(fmt.Sprintf("%s End Live:%d %s\n\n", time.Now(), carClassID, string(jsonResult)))

	return string(jsonResult), nil
}

// provisionalTable calculate change between current and predicted championship tables for the Windows overlay
func provisionalTable(currentStandings, predictedStandings standings.ChampionshipStandings) []live.PredictedStanding {
	mergedStandings := make(map[model.CustID]live.PredictedStanding)

	for _, entry := range currentStandings.Table {
		mergedStandings[entry.CustID] = live.PredictedStanding{
			CurrentPosition: int(entry.Position),
			CustID:          int(entry.CustID),
			DriverName:      entry.DriverName,
			CurrentPoints:   int(entry.DroppedRoundPoints),
		}
	}

	for _, entry := range predictedStandings.Table {
		if _, ok := mergedStandings[entry.CustID]; ok {
			ls := mergedStandings[entry.CustID]

			ls.PredictedPoints = int(entry.DroppedRoundPoints)
			ls.PredictedPosition = int(entry.Position)

			mergedStandings[entry.CustID] = ls
		} else {
			mergedStandings[entry.CustID] = live.PredictedStanding{
				PredictedPosition: int(entry.Position),
				CustID:            int(entry.CustID),
				DriverName:        entry.DriverName,
				PredictedPoints:   int(entry.DroppedRoundPoints),
			}
		}
	}

	predictedResult := make([]live.PredictedStanding, 0, len(mergedStandings))

	for custID := range mergedStandings {
		ls := mergedStandings[custID]
		ls.Change = ls.CurrentPosition - ls.PredictedPosition
		predictedResult = append(predictedResult, ls)
	}

	sort.SliceStable(predictedResult, func(i, j int) bool {
		return predictedResult[i].PredictedPosition < predictedResult[j].PredictedPosition
	})

	return predictedResult
}

func buildResults(carClassID model.CarClassID, liveResults []live.CurrentPosition) []results.Results {
	res := make([]results.Results, 0, len(liveResults))

	for _, lr := range liveResults {
		res = append(res, results.Results{
			CustID:                lr.CustID,
			FinishPositionInClass: lr.FinishPositionInClass,
			LapsComplete:          lr.LapsComplete,
			CarID:                 lr.CarID,
			CarClassID:            int(carClassID),
		})
	}

	return res
}

*/
