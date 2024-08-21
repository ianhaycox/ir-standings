// Package predictor works out provisional standings based on the current race positions
package predictor

import (
	"sort"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/ianhaycox/ir-standings/model/telemetry"
)

type Predictor struct {
	previous          *championship.Championship
	previousStandings standings.ChampionshipStandings
	points            points.PointsStructure
	pastResults       []results.Result
	countBestOf       int
}

func NewPredictor(pointsPerSplit points.PointsPerSplit, pastResults []results.Result, countBestOf int) Predictor {
	return Predictor{
		points:      points.NewPointsStructure(pointsPerSplit),
		pastResults: pastResults,
		countBestOf: countBestOf,
	}
}

// Live championship positions
//
// {CarClassID: 84, ShortName: "GTP", Name: "Nissan GTP ZX-T", CarsInClass: []results.CarsInClass{{CarID: 77}}},
// {CarClassID: 83, ShortName: "GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},
func (p *Predictor) Live(td *telemetry.TelemetryData) live.PredictedStandings {
	const (
		carClassID = 84 // TODO all classes as map
	)

	ps := live.PredictedStandings{
		Status:      td.Status,
		TrackName:   td.TrackName,
		CountBestOf: p.countBestOf,
	}

	seriesID := model.SeriesID(td.SeriesID)

	if p.previous == nil {
		p.previous = championship.NewChampionship(seriesID, nil, p.points, p.countBestOf)
		p.previous.LoadRaceData(p.pastResults)
		p.previousStandings = p.previous.Standings(carClassID)
	}

	liveResults := make([]results.Result, 0, len(p.pastResults))
	liveResults = append(liveResults, p.pastResults...)

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
				Results:        buildResults(&td.Cars),
			},
		},
	})

	predicted := championship.NewChampionship(seriesID, nil, p.points, p.countBestOf)

	predicted.LoadRaceData(liveResults)
	predicted.SetCarClasses(p.previous.CarClasses())

	predictedStandings := predicted.Standings(carClassID)

	carClasses := predicted.CarClasses()

	ps.Standings = make(map[model.CarClassID]live.Standing)
	ps.Standings[84] = live.Standing{
		SoFByCarClass:           td.SofByCarClass()[84],
		CarClassID:              carClassID,
		CarClassName:            carClasses.ClassNames()[84].Name(),
		ClassLeaderLapsComplete: model.LapsComplete(td.LeaderLapsComplete(carClassID)),
		Items:                   p.provisionalTable(predictedStandings),
	}

	return ps
}

// provisionalTable calculate change between current and predicted championship tables for the Windows overlay
func (p *Predictor) provisionalTable(predictedStandings standings.ChampionshipStandings) []live.PredictedStanding {
	mergedStandings := make(map[model.CustID]live.PredictedStanding)

	for _, entry := range p.previousStandings.Table {
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
				PredictedPosition: entry.Position,
				CustID:            entry.CustID,
				DriverName:        entry.DriverName,
				PredictedPoints:   entry.DroppedRoundPoints,
				CarNames:          entry.CarNames,
			}
		}
	}

	predictedResult := make([]live.PredictedStanding, 0, len(mergedStandings))

	for custID := range mergedStandings {
		ls := mergedStandings[custID]
		ls.Change = int(ls.CurrentPosition - ls.PredictedPosition)
		predictedResult = append(predictedResult, ls)
	}

	sort.SliceStable(predictedResult, func(i, j int) bool {
		return predictedResult[i].PredictedPosition < predictedResult[j].PredictedPosition
	})

	return predictedResult
}

func buildResults(cars *telemetry.CarsInfo) []results.Results {
	res := make([]results.Results, 0)

	for _, car := range cars {
		if !car.IsRacing() {
			continue
		}

		res = append(res, results.Results{
			CustID:                car.CustID,
			FinishPositionInClass: car.RacePositionInClass,
			LapsComplete:          car.LapsComplete,
			CarID:                 car.CarID,
			CarClassID:            car.CarClassID,
			CarName:               car.CarName,
		})
	}

	return res
}
