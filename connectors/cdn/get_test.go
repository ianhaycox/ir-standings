package cdn

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/test/readers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testData struct {
	Message string `json:"message,omitempty"`
}

func TestResult(t *testing.T) {
	t.Run("Returns data from the CDN", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		var v testData
		vr := testData{Message: "123456"}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)
		mockAPI.EXPECT().Decode(&v, []byte(`{"message":"123456"}`), "application/json").Return(nil).SetArg(0, vr).SetArg(0, vr)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader(`{"message":"123456"}`), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.NoError(t, err)
		assert.Equal(t, "123456", data.Message)
	})
}

func TestResultErrors(t *testing.T) {
	t.Run("Should return error if PrepareRequest fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("prepare failed"))

		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.ErrorContains(t, err, "prepare failed")
	})

	t.Run("Should return error if CallAPI fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		mockAPI.EXPECT().PrepareRequest(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(request, nil)
		mockAPI.EXPECT().CallAPI(request).Return(nil, fmt.Errorf("call failed"))

		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.ErrorContains(t, err, "call failed")
	})

	t.Run("Should return error if can't read response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		mockAPI.EXPECT().PrepareRequest(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       readers.NewReadCloserErrorWrapper(func(p []byte) (n int, err error) { return 0, fmt.Errorf("read error") }, func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.ErrorContains(t, err, "read error")
	})

	t.Run("Should return error if not 200 response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		response := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("body"), func() error { return nil }),
		}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)
		mockAPI.EXPECT().ReportError(&AWSErrorResponse{}, response, []byte("body")).Return(fmt.Errorf("forbidden"))
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.ErrorContains(t, err, "forbidden")
	})

	t.Run("Should return error if can not decode response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("body"), func() error { return nil }),
		}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		mockAPI.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("decode failed"))

		c := NewCDNService(mockAPI)

		var data testData

		err := c.Get(ctx, "https://cdn.com/results", &data)
		assert.ErrorContains(t, err, "decode failed")
	})
}
