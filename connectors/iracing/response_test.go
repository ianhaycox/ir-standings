package iracing

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method))
}

func TestErrorResponse(t *testing.T) {
	t.Run("without err", func(t *testing.T) {
		assert.Equal(t, `{"code":1,"message":"message"}`, ErrorResponse(1, "message", nil))
	})

	t.Run("with err", func(t *testing.T) {
		assert.Equal(t, `{"code":1,"message":"message, err: an error"}`, ErrorResponse(1, "message", fmt.Errorf("an error")))
	})
}
