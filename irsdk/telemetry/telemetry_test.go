package telemetry

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/ianhaycox/ir-standings/irsdk"
)

func TestTelemetry(t *testing.T) {
	t.Skip()

	reader, err := os.Open("/tmp/test.ibt")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Init irSDK Linux(other)")

	sdk := irsdk.Init(reader)

	online := true
	for {
		sdk.WaitForData(100 * time.Millisecond)

		speed, err := sdk.GetVar("Speed")
		if err != nil {
			t.Fail()
		}

		log.Println(speed)

		if sdk.IsConnected() {
			time.Sleep(500 * time.Millisecond)
			if !online {
				log.Println("iRacing connected!")
			}
			online = true
		} else {
			time.Sleep(5 * time.Second)
			if online {
				log.Println("Waiting for iRacing connection...")
			}
			online = false
		}
	}
}
