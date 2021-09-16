package sources

import (
	"fmt"
	"io"
	"net/http"
)

type SourceHTTP struct {
	URL string
}

func (s *SourceHTTP) Fetch() (io.ReadCloser, error) {
	resp, err := http.Get(s.URL)
	if err != nil {
		return nil, fmt.Errorf("http source failed to fetch: %s", err)
	}

	return resp.Body, nil
}
