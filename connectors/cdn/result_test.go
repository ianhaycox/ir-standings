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

func TestResult(t *testing.T) {
	t.Run("Returns a result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		c := NewCDNService(mockAPI)

		res, err := c.GetResult(ctx, "https://cdn.com/results")
		assert.NoError(t, err)
		assert.Equal(t, "result", res)
	})
}

func TestResultErrors(t *testing.T) {
	t.Run("Should return error if PrepareRequest fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(nil, fmt.Errorf("prepare failed"))

		c := NewCDNService(mockAPI)

		res, err := c.GetResult(ctx, "https://cdn.com/results")
		assert.ErrorContains(t, err, "prepare failed")
		assert.Equal(t, "", res)
	})

	t.Run("Should return error if CallAPI fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)

		mockAPI.EXPECT().CallAPI(request).Return(nil, fmt.Errorf("call failed"))
		c := NewCDNService(mockAPI)

		res, err := c.GetResult(ctx, "https://cdn.com/results")
		assert.ErrorContains(t, err, "call failed")
		assert.Equal(t, "", res)
	})

	t.Run("Should return error if can't read response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		mockAPI := api.NewMockAPIClientInterface(ctrl)

		request := &http.Request{}
		mockAPI.EXPECT().PrepareRequest(ctx, "https://cdn.com/results", "GET", nil, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       readers.NewReadCloserErrorWrapper(func(p []byte) (n int, err error) { return 0, fmt.Errorf("read error") }, func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		c := NewCDNService(mockAPI)

		res, err := c.GetResult(ctx, "https://cdn.com/results")
		assert.ErrorContains(t, err, "read error")
		assert.Equal(t, "", res)
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

		res, err := c.GetResult(ctx, "https://cdn.com/results")
		assert.ErrorContains(t, err, "forbidden")
		assert.Equal(t, "", res)
	})
}
