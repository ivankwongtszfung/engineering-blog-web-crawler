package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var serverIndexResponse = []byte("hello world\n")

func readTestHtml(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func newUnstartedTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(serverIndexResponse)
	})

	mux.HandleFunc("/uber.html", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("no-content-type") != "" {
			w.Header()["Content-Type"] = nil
		} else {
			w.Header().Set("Content-Type", "text/html")
		}
		filename := "./testcases/uber_blog.html"
		bytes, err := readTestHtml(filename)
		if err != nil {
			log.Fatalf("Error when reading %s, %s", filename, err)
		}
		w.Write(bytes)
	})

	return httptest.NewUnstartedServer(mux)
}

func newTestServer() *httptest.Server {
	srv := newUnstartedTestServer()
	srv.Start()
	return srv
}

func TestScrapeUberBlogsToCsv(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()
}
