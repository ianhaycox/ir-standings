package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ianhaycox/ir-standings/test/readers"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method))
}

func TestBodyClose(t *testing.T) {
	t.Parallel()

	t.Run("bodyclose", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		res := w.Result()
		defer BodyClose(res)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, []byte("GET"), data)
	})

	t.Run("bodyclose with nil response", func(t *testing.T) {
		BodyClose(nil)
	})

	t.Run("bodyclose with error response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaders := readers.NewMockReadCloserError(ctrl)
		mockReaders.EXPECT().Close().Return(fmt.Errorf("close error"))

		response := &http.Response{
			Body: mockReaders,
		}

		BodyClose(response)
	})
}
