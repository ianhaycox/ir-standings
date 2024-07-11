// Package iracing declares common api functions
package iracing

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ItemRequired    = 1
	NoPermission    = 2
	JSONError       = 3
	DatabaseError   = 4
	ValidationError = 5
	APIError        = 6
	XMLError        = 7
	CommsError      = 8
)

type APIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BodyClose(response *http.Response) {
	if response != nil && response.Body != nil {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}
}

func ErrorResponse(code int, message string, err error) string {
	apiError := APIErrorResponse{
		Code: code,
	}

	if err != nil {
		apiError.Message = fmt.Sprintf("%s, err: %v", message, err)
	} else {
		apiError.Message = message
	}

	response, err := json.Marshal(apiError)
	if err != nil {
		response = []byte("Unknown error")
	}

	return string(response)
}
