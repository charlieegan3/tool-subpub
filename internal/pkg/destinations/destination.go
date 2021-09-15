package destinations

import "io"

type Destination interface {
	Put(io.ReadCloser) error
}
