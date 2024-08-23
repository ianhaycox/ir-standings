// Package predictor works out provisional standings based on the current race positions
package predictor

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/ianhaycox/ir-standings/model"
	"github.com/ianhaycox/ir-standings/model/championship"
	"github.com/ianhaycox/ir-standings/model/championship/points"
	"github.com/ianhaycox/ir-standings/model/championship/standings"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/live"
	"github.com/ianhaycox/ir-standings/model/telemetry"
	"github.com/ianhaycox/ir-standings/test/devmode"
)

type Predictor struct {
	previous          *championship.Championship
	previousStandings map[model.CarClassID]standings.ChampionshipStandings
	points            points.PointsStructure
	pastResults       []results.Result
	countBestOf       int
	td                *telemetry.TelemetryData
}

func NewPredictor(pointsPerSplit points.PointsPerSplit, pastResults []results.Result, td *telemetry.TelemetryData, countBestOf int) *Predictor {
	if td != nil && devmode.IsDevMode() {
		b, _ := json.MarshalIndent(td, "", "  ")
		_ = os.WriteFile("/tmp/telemetry.json", b, 0600) //nolint:gosec,mnd // ok
	}

	return &Predictor{
		points:            points.NewPointsStructure(pointsPerSplit),
		pastResults:       pastResults,
		countBestOf:       countBestOf,
		td:                td,
		previousStandings: make(map[model.CarClassID]standings.ChampionshipStandings),
	}
}

// Live championship positions
//
// {CarClassID: 84, ShortName: "GTP", Name: "Nissan GTP ZX-T", CarsInClass: []results.CarsInClass{{CarID: 77}}},
// {CarClassID: 83, ShortName: "GTO", Name: "Audi 90 GTO", CarsInClass: []results.CarsInClass{{CarID: 76}}},
func (p *Predictor) Live() live.PredictedStandings {
	const (
		carClassIDNissan = 84
		carClassIDAudi   = 83
	)

	ps := live.PredictedStandings{
		Status:      p.td.Status,
		TrackName:   p.td.TrackName,
		CountBestOf: p.countBestOf,
	}

	seriesID := model.SeriesID(p.td.SeriesID)

	if p.previous == nil {
		p.previous = championship.NewChampionship(seriesID, nil, p.points, p.countBestOf)
		p.previous.LoadRaceData(p.pastResults)
	}

	ps.Standings = make(map[model.CarClassID]live.Standing)

	ps.Standings[carClassIDNissan] = p.liveStandings(seriesID, carClassIDNissan)
	ps.Standings[carClassIDAudi] = p.liveStandings(seriesID, carClassIDAudi)

	return ps
}

func (p *Predictor) liveStandings(seriesID model.SeriesID, carClassID model.CarClassID) live.Standing {
	if _, ok := p.previousStandings[carClassID]; !ok {
		p.previousStandings[carClassID] = p.previous.Standings(carClassID)
	}

	liveResults := make([]results.Result, 0, len(p.pastResults))
	liveResults = append(liveResults, p.pastResults...)

	liveResults = append(liveResults, results.Result{
		SessionID:     p.td.SessionID,
		SubsessionID:  p.td.SubsessionID,
		SeriesID:      p.td.SeriesID,
		SessionSplits: []results.SessionSplits{{SubsessionID: p.td.SubsessionID}},
		StartTime:     time.Now().UTC(),
		Track:         results.ResultTrack{TrackID: p.td.TrackID, TrackName: p.td.TrackName},
		CarClasses:    p.td.CarClasses(),
		SessionResults: []results.SessionResults{
			{
				SimsessionName: "RACE",
				Results:        buildResults(&p.td.Cars),
			},
		},
	})

	predicted := championship.NewChampionship(seriesID, nil, p.points, p.countBestOf)

	predicted.LoadRaceData(liveResults)
	predicted.SetCarClasses(predicted.CarClasses())

	predictedStandings := predicted.Standings(carClassID)

	carClasses := predicted.CarClasses()

	return live.Standing{
		SoFByCarClass:           p.td.SofByCarClass()[int(carClassID)],
		CarClassID:              carClassID,
		CarClassName:            carClasses.ClassNames()[carClassID].Name(),
		ClassLeaderLapsComplete: model.LapsComplete(p.td.LeaderLapsComplete(int(carClassID))),
		Items:                   p.provisionalTable(predictedStandings, carClassID),
	}
}

// provisionalTable calculate change between current and predicted championship tables for the Windows overlay
func (p *Predictor) provisionalTable(predictedStandings standings.ChampionshipStandings, carClassID model.CarClassID) []live.PredictedStanding {
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
	for i := range p.td.Cars {
		carNums[model.CustID(p.td.Cars[i].CustID)] = p.td.Cars[i].CarNumber
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
			DisplayName:           car.DriverName,
			NewiRating:            car.IRating,
		})
	}

	return res
}
