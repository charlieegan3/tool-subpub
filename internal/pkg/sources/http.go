package sources

import "io"

type SourceHTTP struct {
	URL string
}

func (s *SourceHTTP) Fetch() (io.ReadCloser, error) {
	return nil, nil
}
