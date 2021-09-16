package destinations

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestDestinationBucketObject(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "sub-pub-destination-test-")
	if err != nil {
		t.Fatalf("unexpected error making temp file: %s", err)
	}
	defer os.Remove(file.Name())

	fmt.Println(file.Name())

	destination := DestinationBucketObject{
		Path: fmt.Sprintf("file://%s", file.Name()),
	}

	err = destination.Put(ioutil.NopCloser(strings.NewReader("foobar")))
	if err != nil {
		t.Fatalf("unexpected error putting data: %s", err)
	}

	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("unexpected error checking data: %s", err)
	}

	td.Cmp(t, string(bytes), "foobar")
}
