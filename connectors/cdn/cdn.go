//go:generate mockgen -package cdn -destination cdn_mock.go -source cdn.go

// Package cdn gets data from the iRacing Content Delivery Network
package cdn

import (
	"context"

	"github.com/ianhaycox/ir-standings/connectors/api"
)

type AWSErrorResponse struct {
	Code       string `xml:"Code"`
	Message    string `xml:"Message"`
	Expires    string `xml:"Expires"`
	ServerTime string `xml:"ServerTime"`
	RequestID  string `xml:"RequestId"`
	HostID     string `xml:"HostId"`
}

type CDNService struct {
	client api.APIClientInterface
}

type CDNAPI interface {
	GetResult(ctx context.Context, link string) (string, error)
}

// NewCDNService gets iRacing results, etc. from AWS
func NewCDNService(client api.APIClientInterface) *CDNService {
	return &CDNService{
		client: client,
	}
}
