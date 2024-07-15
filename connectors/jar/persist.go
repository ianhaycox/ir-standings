package cookiejar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const mode = 0600

type Store struct {
	filename string
}

type Persister interface {
	Read() (map[string][]*http.Cookie, error)
	Write(map[string][]*http.Cookie) error
}

func NewStore(filename string) Persister {
	return &Store{
		filename: filename,
	}
}

func (s *Store) Read() (map[string][]*http.Cookie, error) {
	cookies := make(map[string][]*http.Cookie)

	b, err := os.ReadFile(s.filename)
	if err != nil {
		return cookies, err
	}

	err = json.Unmarshal(b, &cookies)
	if err != nil {
		return make(map[string][]*http.Cookie), err
	}

	return cookies, nil
}

func (s *Store) Write(cookies map[string][]*http.Cookie) error {
	b, err := json.Marshal(cookies)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filename, b, mode)
	if err != nil {
		return fmt.Errorf("can not close cookie jar:%s, err:%w", s.filename, err)
	}

	return nil
}
