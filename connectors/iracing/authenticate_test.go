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

func TestAuthenticate(t *testing.T) {
	t.Run("Returns no error if a successful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}, nil)

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		var ok AuthenticationGoodResponse
		mockAPI.EXPECT().Decode(&ok, []byte("result"), "application/json").Return(nil)

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.NoError(t, err)
	})

	t.Run("Returns an error if an unsuccessful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}, nil)

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader(`{"authcode":0}`), func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		var ok AuthenticationGoodResponse
		mockAPI.EXPECT().Decode(&ok, []byte(`{"authcode":0}`), "application/json").Return(fmt.Errorf("no auth token"))

		bad := AuthenticationBadResponse{AuthCode: 0}
		mockAPI.EXPECT().Decode(&bad, []byte(`{"authcode":0}`), "application/json").Return(nil)

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "failed to authenticate")
	})
}

func TestAuthenticateErrors(t *testing.T) {
	t.Run("Should return error if can not get credentials fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(nil, fmt.Errorf("no creds"))

		ir := NewIracingService(nil, mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "no creds")
	})

	t.Run("Returns error if PrepareRequest fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{}, nil)

		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{}).Return(nil, fmt.Errorf("prepare failed"))

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "prepare failed")
	})

	t.Run("Returns error if callAPI fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{}, nil)

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{}).Return(request, nil)

		mockAPI.EXPECT().CallAPI(request).Return(nil, fmt.Errorf("callapi failed"))

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "callapi failed")
	})

	t.Run("Returns error if can not read response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{}, nil)

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{}).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       readers.NewReadCloserErrorWrapper(func(p []byte) (n int, err error) { return 0, fmt.Errorf("read error") }, func() error { return nil }),
		}
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "read error")
	})

	t.Run("Returns error if not a 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockAuth := api.NewMockAuthenticator(ctrl)
		mockAuth.EXPECT().Credentials().Return(&api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}, nil)

		request := &http.Request{}
		mockAPI := api.NewMockAPIClientInterface(ctrl)
		mockAPI.EXPECT().PrepareRequest(ctx, "https://members-ng.iracing.com/auth", "POST", url.Values{}, &api.Credentials{Email: "test@example.com", EncodedPassword: "1234"}).Return(request, nil)

		response := &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       ioutils.NewReadCloserWrapper(strings.NewReader("result"), func() error { return nil }),
		}
		var apiError APIErrorResponse
		mockAPI.EXPECT().CallAPI(request).Return(response, nil)
		mockAPI.EXPECT().ReportError(&apiError, response, []byte("result")).Return(fmt.Errorf("403"))

		ir := NewIracingService(NewIracingDataService(mockAPI), mockAuth)

		err := ir.Authenticate(ctx)
		assert.ErrorContains(t, err, "403")
	})
}
