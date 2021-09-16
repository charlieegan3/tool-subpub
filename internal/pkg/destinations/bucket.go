package destinations

import (
	"context"
	"fmt"
	"io"
	"strings"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

type DestinationBucketObject struct {
	Path string
}

func (d *DestinationBucketObject) Put(input io.ReadCloser) error {
	defer input.Close()

	pathParts := strings.Split(d.Path, "/")
	if len(pathParts) < 2 {
		return fmt.Errorf("invalid path, cannot determine file name and dir")
	}

	file := pathParts[len(pathParts)-1]
	dir := strings.Join(pathParts[0:len(pathParts)-1], "/")

	b, err := blob.OpenBucket(context.TODO(), dir)
	if err != nil {
		return err
	}

	w, err := b.NewWriter(context.TODO(), file, nil)
	if err != nil {
		return fmt.Errorf("failed to open bucket: %s", err)
	}
	if _, err := io.Copy(w, input); err != nil {
		return fmt.Errorf("failed to copy data to path: %s", err)
	}
	closeErr := w.Close()
	if closeErr != nil {
		return fmt.Errorf("failed to close bucket writer: %s", err)
	}

	closeErr = b.Close()
	if closeErr != nil {
		return fmt.Errorf("failed to close bucket: %s", err)
	}

	return nil
}
