package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
}
