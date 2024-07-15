//go:generate mockgen -package iracing -destination iracing_mock.go -source iracing.go

// Package iracing API
package iracing

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

const (
	Endpoint string = "https://members-ng.iracing.com" // No trailing /
)

var (
	CookiesFile = filepath.Join(os.TempDir(), "ir-standings-cookies")
)

type IracingAPI interface {
	Authenticate(ctx context.Context)
	ResultLink(ctx context.Context, subsessionID string) (ResultsLink, error)
}

type IracingDataAPI interface {
	Get(ctx context.Context, v interface{}, path string, queryParams url.Values) error // GET from the data endpoint
	Client() api.APIClientInterface
}

type IracingDataService struct {
	client api.APIClientInterface
}

func (ids *IracingDataService) Client() api.APIClientInterface {
	return ids.client
}

type IracingService struct {
	data IracingDataAPI
	auth api.Authenticator
}

type APIErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Note    string `json:"note,omitempty"`
}

func NewIracingDataService(client api.APIClientInterface) *IracingDataService {
	return &IracingDataService{client}
}

// NewIracingService https://members-ng.iracing.com/data/doc
func NewIracingService(data IracingDataAPI, auth api.Authenticator) *IracingService {
	return &IracingService{
		data: data,
		auth: auth,
	}
}
