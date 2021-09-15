package sources

import "io"

type Source interface {
	Fetch() (io.ReadCloser, error)
}
