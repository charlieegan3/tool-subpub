package sources

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestSourceHTTP(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td.Cmp(t, r.Method, "GET")

		w.Write([]byte("response"))
	}))

	source := SourceHTTP{URL: testServer.URL}

	reader, err := source.Fetch()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	td.Cmp(t, string(output), "response")
}
