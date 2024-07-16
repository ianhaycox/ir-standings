package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
)

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
	case "result":
		var obj results.Result

		err := json.Unmarshal(buf, &obj)
		if err != nil {
			panic(err)
		}

		for i := range obj.SessionResults {
			for j := range obj.SessionResults[i].Results {
				obj.SessionResults[i].Results[j].DisplayName = redact(obj.SessionResults[i].Results[j].DisplayName)
			}
		}

		readactedBuf, err = json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}

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
	}

	err = os.WriteFile(os.Args[1], readactedBuf, 0600) //nolint:mnd // ok
	if err != nil {
		panic(err)
	}
}

func redact(s string) string {
	var res string

	for i := 0; i < len(s); i++ {
		if i%2 == 0 {
			res += s[i : i+1]
		} else {
			res += "*"
		}
	}

	return res
}
