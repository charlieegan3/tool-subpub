package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/charlieegan3/tool-subpub/pkg/api"
)

func BuildGetHandler(targets map[string]api.Target) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error
		vars := mux.Vars(request)

		target, ok := vars["target"]
		if !ok || target == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		loadedTarget, ok := targets[target]
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		// make a request to loadedTarget.URL and get the content

		client := http.Client{}
		resp, err := client.Get(loadedTarget.URL)
		if err != nil {
			writer.WriteHeader(resp.StatusCode)
			writer.Write([]byte(err.Error()))
			return
		}
		if resp.StatusCode >= 400 {
			writer.WriteHeader(resp.StatusCode)
			writer.Write([]byte(err.Error()))
			return
		}

		for k, v := range resp.Header {
			if k == "Content-Length" {
				continue
			}
			writer.Header().Set(k, v[0])
		}

		reader := resp.Body
		for _, s := range loadedTarget.Substitutions {
			reader, err = s.Run(reader)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
				return
			}
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
	}
}
