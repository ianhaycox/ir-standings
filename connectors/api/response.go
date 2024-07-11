// Package api declares common api functions
package api

import (
	"net/http"
)

type APIErrorResponse struct {
	Error int `json:"error"`
}

func BodyClose(response *http.Response) {
	if response != nil && response.Body != nil {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}
}
