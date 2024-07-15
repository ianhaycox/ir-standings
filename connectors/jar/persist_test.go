package cookiejar

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistence(t *testing.T) {
	t.Run("Reads an existing cookiejar OK", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "example")
		assert.NoError(t, err)

		defer os.RemoveAll(tempDir)

		cookiesFile := filepath.Join(tempDir, "cookies")
		err = os.WriteFile(cookiesFile, []byte(`{"host":[{"Name":"name"}]}`), mode)
		require.NoError(t, err)

		s := NewStore(cookiesFile)

		cookies, err := s.Read()
		assert.NoError(t, err)
		assert.Equal(t, "name", cookies["host"][0].Name)
	})

	t.Run("Writes cookes to a cookiejar", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "example")
		assert.NoError(t, err)

		defer os.RemoveAll(tempDir)

		cookiesFile := filepath.Join(tempDir, "cookies")

		s := NewStore(cookiesFile)

		cookies := make(map[string][]*http.Cookie)
		cookies["host"] = append(cookies["host"], &http.Cookie{Name: "name"})

		err = s.Write(cookies)
		assert.NoError(t, err)

		b, err := os.ReadFile(cookiesFile)
		assert.NoError(t, err)
		assert.Contains(t, string(b), `"Name":"name"`)
	})
}

func TestPersistenceErrors(t *testing.T) {
	t.Run("Returns an error if the cookiejar is not readable", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "example")
		assert.NoError(t, err)

		defer os.RemoveAll(tempDir)

		cookiesFile := filepath.Join(tempDir, "cookies")

		s := NewStore(cookiesFile)

		cookies, err := s.Read()
		assert.Error(t, err)
		assert.Equal(t, make(map[string][]*http.Cookie), cookies)
	})

	t.Run("Returns an error if the cookiejar contents are invalid", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "example")
		assert.NoError(t, err)

		defer os.RemoveAll(tempDir)

		cookiesFile := filepath.Join(tempDir, "cookies")
		err = os.WriteFile(cookiesFile, []byte(`junk`), mode)
		require.NoError(t, err)

		s := NewStore(cookiesFile)

		cookies, err := s.Read()
		assert.Error(t, err)
		assert.Equal(t, make(map[string][]*http.Cookie), cookies)
	})

	t.Run("Returns an error if the cookiejar is not writeable", func(t *testing.T) {
		s := NewStore("/does/not/exist")

		err := s.Write(make(map[string][]*http.Cookie))
		assert.Error(t, err)
	})
}
