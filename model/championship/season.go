package championship

import (
	"fmt"
	"sort"

	"github.com/ianhaycox/ir-standings/model/data/results"
)

type Season struct {
	events     []Event              // In order of event
	carClasses []results.CarClasses // Competing cars
	splits     []SplitResult        // Splits index + 1
}

type SplitResult struct {
	results          map[SubsessionID]map[CustID]*Result
	pointsByCarClass map[CarClassID]map[CustID]SeasonPoints
}

type SeasonPoints struct {
	custID              CustID
	displayName         string
	bestOf              int
	finishingPositions  map[SubsessionID]int
	championshipsPoints map[SubsessionID]int
}

func (s *Season) CalculatePositionsByCarClass() {
	for _, splitResult := range s.splits {
		for _, result := range splitResult.results {
			for custID, sessionResult := range result {
				carClassID := sessionResult.CarClassID

				if _, ok := splitResult.pointsByCarClass[carClassID]; !ok {
					splitResult.pointsByCarClass[carClassID] = make(map[CustID]SeasonPoints)
				}

				if _, ok := splitResult.pointsByCarClass[carClassID][custID]; !ok {
					seasonPoints := SeasonPoints{
						custID:              custID,
						displayName:         sessionResult.DisplayName,
						finishingPositions:  make(map[SubsessionID]int),
						championshipsPoints: make(map[SubsessionID]int),
					}

					splitResult.pointsByCarClass[carClassID][custID] = seasonPoints
				}

				splitResult.pointsByCarClass[carClassID][custID].finishingPositions[sessionResult.SubsessionID] = sessionResult.FinishPositionInClass
			}
		}
	}
}

func (s *Season) CalculateChampionshipPoints() {
	var points = []int{25, 22, 20, 18, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	for _, splitResult := range s.splits {
		winnerLapsComplete := make(map[SubsessionID]map[CarClassID]int)

		for subsessionID, result := range splitResult.results {
			for _, sessionResult := range result {
				carClassID := sessionResult.CarClassID

				if len(winnerLapsComplete[subsessionID]) == 0 {
					winnerLapsComplete[subsessionID] = make(map[CarClassID]int)
				}

				if sessionResult.LapsComplete > winnerLapsComplete[subsessionID][carClassID] {
					winnerLapsComplete[subsessionID][carClassID] = sessionResult.LapsComplete
				}
			}
		}

		for carClassID, car := range splitResult.pointsByCarClass {
			for custID := range car {
				for subsessionID := range car[custID].finishingPositions {
					if splitResult.results[subsessionID][custID].IsClassified(winnerLapsComplete[subsessionID][carClassID]) {
						if len(points) > splitResult.pointsByCarClass[carClassID][custID].finishingPositions[subsessionID] {
							splitResult.pointsByCarClass[carClassID][custID].championshipsPoints[subsessionID] =
								points[splitResult.pointsByCarClass[carClassID][custID].finishingPositions[subsessionID]]
						}
					}
				}

				all := make([]int, 0)
				for _, champPoints := range splitResult.pointsByCarClass[carClassID][custID].championshipsPoints {
					all = append(all, champPoints)
				}

				sort.SliceStable(all, func(i, j int) bool { return all[i] > all[j] })

				bestOf := 0

				for i := 0; i < 9 && i < len(all); i++ {
					bestOf += all[i]
				}

				x := splitResult.pointsByCarClass[carClassID][custID]
				x.bestOf = bestOf
				splitResult.pointsByCarClass[carClassID][custID] = x
			}
		}
	}
}

type tab struct {
	DisplayName string
	BestOf      int
}

func (s *Season) PrintTable() {
	for splitNum, splitResult := range s.splits {
		fmt.Printf("Split:%d\n", splitNum)

		for carClassID, car := range splitResult.pointsByCarClass {
			fmt.Printf("  Car:%d\n", carClassID)

			out := make([]tab, 0)

			for custID := range car {
				out = append(out, tab{DisplayName: car[custID].displayName, BestOf: car[custID].bestOf})

				//fmt.Printf("    %s %d\n", car[custID].displayName, car[custID].bestOf)
			}

			sort.SliceStable(out, func(i, j int) bool { return out[i].BestOf > out[j].BestOf })

			for _, l := range out {
				fmt.Printf("    %s %d\n", l.DisplayName, l.BestOf)
			}
		}

	}

}
