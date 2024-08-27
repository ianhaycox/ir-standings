//go:generate mockgen -package iracing -destination iracing_mock.go -source iracing.go

// Package iracing API
package iracing

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/connectors/cdn"
	"github.com/ianhaycox/ir-standings/model/data/cars"
	"github.com/ianhaycox/ir-standings/model/data/results"
	"github.com/ianhaycox/ir-standings/model/data/results/searchseries"
	"github.com/ianhaycox/ir-standings/model/data/seasons"
)

const (
	Endpoint string = "https://members-ng.iracing.com" // No trailing /
)

var (
	CookiesFile = filepath.Join(os.TempDir(), "ir-standings-cookies")
)

type APIErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Note    string `json:"note,omitempty"`
}

type IracingService interface {
	Authenticate(ctx context.Context, email, password string) error
	ResultLink(ctx context.Context, subsessionID int) (*results.ResultLink, error)
	SearchSeriesResults(ctx context.Context, seasonYear, seasonQuarter, seriesID int) ([]searchseries.SearchSeriesResult, error)
	SeasonBroadcastResults(ctx context.Context, ssResults []searchseries.SearchSeriesResult) ([]results.Result, error)
	Seasons(ctx context.Context) ([]seasons.Season, error)
	Cars(ctx context.Context) ([]cars.Car, error)
	CarClasses(ctx context.Context) ([]cars.CarClass, error)
}

type IracingAPI struct {
	client api.API
	data   IracingDataService
	auth   api.Authenticator
}

// NewIracingService https://members-ng.iracing.com/data/doc
func NewIracingService(client api.API, data IracingDataService, auth api.Authenticator) *IracingAPI {
	return &IracingAPI{
		client: client,
		data:   data,
		auth:   auth,
	}
}

type IracingDataService interface {
	Get(ctx context.Context, v interface{}, path string, queryParams url.Values) error // GET from the members /data endpoint
	cdn.CDNAPI                                                                         // GET from the Amazon bucket
}

type IracingDataAPI struct {
	client api.API
	cdn    cdn.CDNAPI
}

func (ids *IracingDataAPI) CDN(ctx context.Context, link string, v interface{}) error {
	return ids.cdn.CDN(ctx, link, v)
}

func NewIracingDataService(client api.API, cdn cdn.CDNAPI) *IracingDataAPI {
	return &IracingDataAPI{
		client: client,
		cdn:    cdn,
	}
}
