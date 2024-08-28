// Package predictor works out provisional standings based on the current race positions
package predictor

import (
	"sort"
	"time"

	"github.com/ianhaycox/ir-standings/irsdk/telemetry"
	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship"
	"github.com/ianhaycox/ir-standings/model/championship/car"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
)

type Predictor struct {
	previous          *championship.Championship
	previousStandings map[model.CarClassID]standings.ChampionshipStandings
	points            points.PointsStructure
	countBestOf       int
	carClasses        car.CarClasses
}

func NewPredictor(pointsPerSplit points.PointsPerSplit, countBestOf int, carClasses car.CarClasses) *Predictor {
	return &Predictor{
		points:            points.NewPointsStructure(pointsPerSplit),
		countBestOf:       countBestOf,
		previousStandings: make(map[model.CarClassID]standings.ChampionshipStandings),
		carClasses:        carClasses,
	}
}

// Live championship positions
//
// {CarClassID: 84, ShortName: "GTP", Name: "Nissan GTP ZX-T", CarsInClass: []results.CarsInClass{{CarID: 77}}},
// {CarClassID: 83, ShortName: "GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},
func (p *Predictor) Live(pastResults []results.Result, td *telemetry.TelemetryData) live.PredictedStandings {
	ps := live.PredictedStandings{
		Status:         td.Status,
		TrackName:      td.TrackName,
		CountBestOf:    p.countBestOf,
		SelfCarClassID: td.SelfCarClassID,
		CarClassIDs:    p.carClasses.CarClassIDs(),
	}

	seriesID := model.SeriesID(td.SeriesID)

	if p.previous == nil {
		p.previous = championship.NewChampionship(seriesID, p.carClasses, nil, p.points, p.countBestOf)
		p.previous.LoadRaceData(pastResults)
	}

	ps.Standings = make(map[model.CarClassID]live.Standing)

	for _, carClassID := range ps.CarClassIDs {
		cci := model.CarClassID(carClassID)
		ps.Standings[cci] = p.liveStandings(seriesID, cci, pastResults, td)
	}

	return ps
}

func (p *Predictor) liveStandings(seriesID model.SeriesID, carClassID model.CarClassID, pastResults []results.Result,
	td *telemetry.TelemetryData) live.Standing {
	if _, ok := p.previousStandings[carClassID]; !ok {
		p.previousStandings[carClassID] = p.previous.Standings(carClassID)
	}

	liveResults := make([]results.Result, 0, len(pastResults))
	liveResults = append(liveResults, pastResults...)

	liveResults = append(liveResults, results.Result{
		SessionID:     td.SessionID,
		SubsessionID:  td.SubsessionID,
		SeriesID:      td.SeriesID,
		SessionSplits: []results.SessionSplits{{SubsessionID: td.SubsessionID}},
		StartTime:     time.Now().UTC(),
		Track:         results.ResultTrack{TrackID: td.TrackID, TrackName: td.TrackName},
		SessionResults: []results.SessionResults{
			{
				SimsessionName: "RACE",
				Results:        p.buildResults(&td.Cars, td.SessionType),
			},
		},
	})

	predicted := championship.NewChampionship(seriesID, p.carClasses, nil, p.points, p.countBestOf)

	predicted.LoadRaceData(liveResults)

	predictedStandings := predicted.Standings(carClassID)

	return live.Standing{
		SoFByCarClass:           td.SofByCarClass()[int(carClassID)],
		CarClassID:              carClassID,
		CarClassName:            p.carClasses.Name(carClassID),
		ClassLeaderLapsComplete: model.LapsComplete(td.LeaderLapsComplete(int(carClassID))),
		Items:                   p.provisionalTable(predictedStandings, carClassID, td),
	}
}

// provisionalTable calculate change between current and predicted championship tables for the Windows overlay
func (p *Predictor) provisionalTable(predictedStandings standings.ChampionshipStandings, carClassID model.CarClassID,
	td *telemetry.TelemetryData) []live.PredictedStanding {
	mergedStandings := make(map[model.CustID]live.PredictedStanding)

	for _, entry := range p.previousStandings[carClassID].Table {
		mergedStandings[entry.CustID] = live.PredictedStanding{
			CurrentPosition: entry.Position,
			CustID:          entry.CustID,
			DriverName:      entry.DriverName,
			CurrentPoints:   entry.DroppedRoundPoints,
			CarNames:        entry.CarNames,
		}
	}

	for _, entry := range predictedStandings.Table {
		if _, ok := mergedStandings[entry.CustID]; ok {
			ls := mergedStandings[entry.CustID]

			ls.PredictedPoints = entry.DroppedRoundPoints
			ls.PredictedPosition = entry.Position

			mergedStandings[entry.CustID] = ls
		} else {
			mergedStandings[entry.CustID] = live.PredictedStanding{
				CurrentPosition:   entry.Position,
				PredictedPosition: entry.Position,
				CustID:            entry.CustID,
				DriverName:        entry.DriverName,
				PredictedPoints:   entry.DroppedRoundPoints,
				CarNames:          entry.CarNames,
			}
		}
	}

	predictedResult := make([]live.PredictedStanding, 0, len(mergedStandings))

	carNums := make(map[model.CustID]string)
	for i := range td.Cars {
		carNums[model.CustID(td.Cars[i].CustID)] = td.Cars[i].CarNumber
	}

	for custID := range mergedStandings {
		ls := mergedStandings[custID]
		if ls.CurrentPosition != 0 {
			ls.Change = int(ls.CurrentPosition - ls.PredictedPosition)
		}

		ls.CarNumber, ls.Driving = carNums[custID]

		predictedResult = append(predictedResult, ls)
	}

	sort.SliceStable(predictedResult, func(i, j int) bool {
		return predictedResult[i].PredictedPosition < predictedResult[j].PredictedPosition
	})

	return predictedResult
}

// Create a fake result for the race based on current positions
func (p *Predictor) buildResults(cars *telemetry.CarsInfo, sessionType string) []results.Results {
	res := make([]results.Results, 0)

	for _, car := range cars {
		if !car.IsRacing() {
			continue
		}

		racePositionInClass := 0
		lapsComplete := 0

		if sessionType == "RACE" {
			racePositionInClass = car.RacePositionInClass
			lapsComplete = car.LapsComplete
		}

		res = append(res, results.Results{
			CustID:                car.CustID,
			FinishPositionInClass: racePositionInClass,
			LapsComplete:          lapsComplete,
			CarID:                 car.CarID,
			CarClassID:            car.CarClassID,
			DisplayName:           car.DriverName,
			NewiRating:            car.IRating,
		})
	}

	return res
}
