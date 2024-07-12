package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/iracing"
	cookiejar "github.com/ianhaycox/ir-standings/connectors/jar"
)

// Get authentication token from iRacing and store in a cookie jar
func main() {
	const (
		numArgs = 2
	)

	flag.Parse()

	if len(flag.Args()) != numArgs {
		log.Fatal("insufficient args")
	}

	cookiesFile := filepath.Join(os.TempDir(), "ir-standings-cookies")

	ctx := context.Background()
	httpClient := http.DefaultClient
	httpClient.Jar = cookiejar.NewCookieJar(cookiesFile)

	cfg := api.NewConfiguration(iracing.Endpoint, httpClient)
	auth := api.NewAuthenticationService(flag.Arg(0), flag.Arg(1))
	client := api.NewAPIClient(cfg)

	ir := iracing.NewIracingService(client, auth)

	err := ir.Authenticate(ctx)
	if err != nil {
		log.Fatal("Can not login: ", err)
	}

	fmt.Println("Logged in")
}
