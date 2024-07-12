package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResult struct {
	ID string `json:"id,omitempty" xml:"id,omitempty"`
}

func TestPrepareRequest(t *testing.T) {
	t.Parallel()

	t.Run("return error if URL path cannot be parsed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("http://a b", nil))

		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodGet, url.Values{}, nil)
		assert.Error(t, err)
	})

	t.Run("return no error if GET request OK", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("http://bookings", nil))

		queryParams := url.Values{}
		queryParams.Add("foo", "bar")
		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodGet, queryParams, nil)
		assert.NoError(t, err)
	})

	t.Run("return no error if POST request OK", func(t *testing.T) {
		type testPost struct {
			Message string `json:"message,omitempty"`
		}

		data := testPost{Message: "test"}
		api := NewAPIClient(NewConfiguration("http://bookings", nil))

		queryParams := url.Values{}
		queryParams.Add("foo", "bar")
		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodPost, queryParams, data)
		assert.NoError(t, err)
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()

	t.Run("decodes JSON successfully", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		result := testResult{}
		err := api.Decode(&result, []byte(`{"id":"test"}`), "application/json")

		assert.NoError(t, err)
		assert.Equal(t, testResult{ID: "test"}, result)
	})

	t.Run("returns error for unsuccessful JSON decode", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		result := testResult{}
		err := api.Decode(&result, []byte(`{"id":test}`), "application/json")

		assert.Error(t, err)
	})

	t.Run("decodes XML successfully", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		result := testResult{}
		err := api.Decode(&result, []byte(`<xml><id>test</id></xml>`), "application/xml")

		assert.NoError(t, err)
		assert.Equal(t, testResult{ID: "test"}, result)
	})

	t.Run("returns error for unsuccessful XML decode", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		result := testResult{}
		err := api.Decode(&result, []byte(`<xml><id>test</notid></xml>`), "application/xml")

		assert.Error(t, err)
	})

	t.Run("returns error for unknown content type", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		result := testResult{}
		err := api.Decode(&result, []byte(``), "unknown")

		assert.Error(t, err)
	})
}

func TestReportError(t *testing.T) {
	t.Run("should return the error from decode if it fails", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		err := api.ReportError(&http.Response{Header: http.Header{"Content-Type": []string{"application/json"}}}, []byte("bad json response"))
		assert.ErrorContains(t, err, "server returned an error")
		assert.ErrorContains(t, err, "bad json")
		assert.ErrorContains(t, err, "invalid character 'b'")
	})

	t.Run("should return error if the api call failed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		err := api.ReportError(&http.Response{StatusCode: http.StatusBadRequest}, []byte(""))
		assert.ErrorContains(t, err, "server returned non-200")
		assert.ErrorContains(t, err, "400, response ''")
	})

	t.Run("should return error from response if the api call failed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("", nil))

		err := api.ReportError(&http.Response{Header: http.Header{"Content-Type": []string{"application/json"}}, StatusCode: http.StatusBadRequest}, []byte(`{"message":"server error"}`))
		assert.ErrorContains(t, err, "server returned non-200 http code: 400")
		assert.ErrorContains(t, err, "server error")
	})
}
