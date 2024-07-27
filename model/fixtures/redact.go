package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
)

//nolint:gocognit // engage brain
func main() {
	if len(os.Args) != 3 { //nolint:mnd // ok
		log.Fatal("need filename and type")
	}

	buf, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var readactedBuf []byte

	switch os.Args[2] {
	case "results":
		var obj []results.Result

		err := json.Unmarshal(buf, &obj)
		if err != nil {
			panic(err)
		}

		for r := range obj {
			for i := range obj[r].SessionResults {
				for j := range obj[r].SessionResults[i].Results {
					obj[r].SessionResults[i].Results[j].DisplayName = redact(obj[r].SessionResults[i].Results[j].DisplayName)
				}
			}
		}

		readactedBuf, err = json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(readactedBuf))

	case "search-results":
		var obj []searchseries.SearchSeriesResult

		err := json.Unmarshal(buf, &obj)
		if err != nil {
			panic(err)
		}

		for i := range obj {
			obj[i].WinnerName = redact(obj[i].WinnerName)
		}

		readactedBuf, err = json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(readactedBuf))

	case "csv":
		csvReader := csv.NewReader(bytes.NewReader(buf))

		records, err := csvReader.ReadAll()
		if err != nil {
			panic(err)
		}

		for i := range records {
			if i == 0 {
				continue
			}

			records[i][2] = redact(records[i][2])
		}

		w := csv.NewWriter(os.Stdout)

		err = w.WriteAll(records)
		if err != nil {
			panic(err)
		}

	default:
		log.Fatal("need filename and type")
	}
}

func redact(s string) string {
	var res []rune

	i := 2
	for _, c := range s {
		if i%3 == 0 {
			res = append(res, rune('*'))
		} else {
			res = append(res, c)
		}

		i++
	}

	return string(res)
}
