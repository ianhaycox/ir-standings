package iracing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/ianhaycox/ir-standings/connectors/api"
	"github.com/ianhaycox/ir-standings/test/readers"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

type obj struct {
	Value int `json:"value,omitempty"`
}

func TestGet(t *testing.T) {
	t.Run("Returns object successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://example.com/data", "GET", url.Values{"test": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		var v obj
		vRet := obj{Value: 9999}
		mockAPI.EXPECT().Decode(&v, []byte("result"), "application/json").Return(nil).SetArg(0, vRet)

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var actual obj
		expected := obj{Value: 9999}

		err := ir.data.Get(ctx, &actual, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetErrors(t *testing.T) {
	t.Run("Returns error if PrepareRequest fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("prepare failed"))

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var v obj

		err := ir.data.Get(ctx, &v, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.ErrorContains(t, err, "prepare failed")
	})

	t.Run("Returns error if callAPI fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://example.com/data", "GET", url.Values{"test": {"12345"}}, nil).Return(request, nil)

		mockAPI.EXPECT().CallAPI(request).Return(nil, fmt.Errorf("callapi failed"))

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var v obj

		err := ir.data.Get(ctx, &v, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.ErrorContains(t, err, "callapi failed")
	})

	t.Run("Returns error if can not read response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://example.com/data", "GET", url.Values{"test": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       readers.NewReadCloserErrorWrapper(func(p []byte) (n int, err error) { return 0, fmt.Errorf("read error") }, func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var v obj

		err := ir.data.Get(ctx, &v, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.ErrorContains(t, err, "read error")
	})

	t.Run("Returns error if not a 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://example.com/data", "GET", url.Values{"test": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		var apiError APIErrorResponse
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		mockAPI.EXPECT().ReportError(&apiError, response, []byte("result")).Return(fmt.Errorf("403"))

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var v obj

		err := ir.data.Get(ctx, &v, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.ErrorContains(t, err, "403")
	})

	t.Run("Returns error if can not parse response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://example.com/data", "GET", url.Values{"test": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		mockAPI.EXPECT().Decode(gomock.Any(), []byte("result"), "application/json").Return(fmt.Errorf("opps"))

		ir := NewIracingService(NewIracingDataService(mockAPI), nil)

		var v obj

		err := ir.data.Get(ctx, &v, "https://example.com/data", url.Values{"test": {"12345"}})
		assert.ErrorContains(t, err, "could not decode body")
	})
}
