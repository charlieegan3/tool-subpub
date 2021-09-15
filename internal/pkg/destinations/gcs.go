package destinations

import "io"

type DestinationGCSBucketObject struct {
	BucketName string
	ObjectName string
}

func (d *DestinationGCSBucketObject) Put(input io.ReadCloser) error {
	return nil
}
