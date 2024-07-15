package iracing

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestResultLink(t *testing.T) {
	t.Run("Returns no error if successfully gets result link", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		var link ResultsLink
		linkResponse := ResultsLink{Link: "https://cdn.com/result/8w8wh8e"}

		queryParams := url.Values{}
		queryParams.Add("subsession_id", "12345")
		mockDataAPI := NewMockIracingDataAPI(ctrl)
		mockDataAPI.EXPECT().Get(ctx, &link, Endpoint+"/data/results/get", queryParams).Return(nil).SetArg(1, linkResponse)

		ir := NewIracingService(mockDataAPI, nil)

		actual, err := ir.ResultLink(ctx, 12345)
		assert.NoError(t, err)
		assert.Equal(t, "https://cdn.com/result/8w8wh8e", actual.Link)
	})
}

func TestResultLinkErrors(t *testing.T) {
	t.Run("Returns error if Get data fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockDataAPI := NewMockIracingDataAPI(ctrl)
		mockDataAPI.EXPECT().Get(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("get failed"))

		ir := NewIracingService(mockDataAPI, nil)

		link, err := ir.ResultLink(ctx, 12345)
		assert.ErrorContains(t, err, "get failed")
		assert.Nil(t, link)
	})
}
