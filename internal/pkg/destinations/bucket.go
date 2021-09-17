package destinations

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
)

type DestinationBucketObject struct {
	Path string
}

func (d *DestinationBucketObject) Put(input io.ReadCloser) error {
	defer input.Close()

	re := regexp.MustCompile(`^(\w+:\/\/)(\/?[^\/]+)(/.*)$`)
	matches := re.FindStringSubmatch(d.Path)

	if len(matches) != 4 {
		return fmt.Errorf("path not of expected format")
	}

	ctx := context.TODO()
	b, err := blob.OpenBucket(ctx, matches[1]+matches[2])
	if err != nil {
		return fmt.Errorf("failed to open bucket: %s", err)
	}

	// handle difference in the way that gcs and file do abs paths
	filePath := matches[3]
	if matches[1] == "gs://" {
		filePath = strings.TrimPrefix(filePath, "/")
	}

	w, err := b.NewWriter(ctx, filePath, nil)
	if err != nil {
		return fmt.Errorf("failed to write to bucket: %s", err)
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
