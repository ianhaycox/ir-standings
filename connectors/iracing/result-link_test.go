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

func TestResultLink(t *testing.T) {
	t.Run("Returns no error if successfully gets result link", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		resultLink := ResultsLink{}
		resultLinkReturned := ResultsLink{Link: "https://cdn.com/result/8w8wh8e"}
		mockAPI.EXPECT().Decode(&resultLink, []byte("result"), "application/json").Return(nil).SetArg(0, resultLinkReturned)

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.NoError(t, err)
		assert.Equal(t, "https://cdn.com/result/8w8wh8e", link.Link)
	})
}

func TestResultLinkErrors(t *testing.T) {
	t.Run("Returns error if PrepareRequest fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(nil, fmt.Errorf("prepare failed"))

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.ErrorContains(t, err, "prepare failed")
		assert.Nil(t, link)
	})

	t.Run("Returns error if callAPI fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(request, nil)

		mockAPI.EXPECT().CallAPI(request).Return(nil, fmt.Errorf("callapi failed"))

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.ErrorContains(t, err, "callapi failed")
		assert.Nil(t, link)
	})

	t.Run("Returns error if can not read response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       readers.NewReadCloserErrorWrapper(func(p []byte) (n int, err error) { return 0, fmt.Errorf("read error") }, func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.ErrorContains(t, err, "read error")
		assert.Nil(t, link)
	})

	t.Run("Returns error if not a 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		var apiError APIErrorResponse
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		mockAPI.EXPECT().ReportError(&apiError, response, []byte("result")).Return(fmt.Errorf("403"))

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.ErrorContains(t, err, "403")
		assert.Nil(t, link)
	})

	t.Run("Returns error if can not parse response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/data/results/get", "GET", url.Values{"subsession_id": {"12345"}}, nil).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		resultLink := ResultsLink{}
		mockAPI.EXPECT().Decode(&resultLink, []byte("result"), "application/json").Return(fmt.Errorf("opps"))

		ir := NewIracingService(mockAPI, nil)

		link, err := ir.ResultLink(ctx, "12345")
		assert.ErrorContains(t, err, "could not decode body")
		assert.Nil(t, link)
	})
}
