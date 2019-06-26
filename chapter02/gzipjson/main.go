package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	source := map[string]string{
		"Hello": "world",
	}

	gzipWriter := gzip.NewWriter(w)
	writer := io.MultiWriter(os.Stdout, gzipWriter)

	// b, _ := json.Marshal(source)
	// writer.Write(b)

	encoder := json.NewEncoder(writer)
	encoder.Encode(source)

	gzipWriter.Flush()
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("port 8080")
	http.ListenAndServe(":8080", nil)
}
