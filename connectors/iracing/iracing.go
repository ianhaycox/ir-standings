//go:generate mockgen -package iracing -destination iracing_mock.go -source iracing.go

// Package iracing API
package iracing

import (
	"context"
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

type IracingService struct {
	client api.APIClientInterface
	auth   api.Authenticator
}

type APIErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// NewIracingService https://members-ng.iracing.com/data/doc
func NewIracingService(client api.APIClientInterface, auth api.Authenticator) *IracingService {
	return &IracingService{
		client: client,
		auth:   auth,
	}
}

type IracingAPI interface {
	Authenticate(ctx context.Context)
	ResultLink(ctx context.Context, subsessionID string) (ResultsLink, error)
}
