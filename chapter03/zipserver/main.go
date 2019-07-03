package main

import (
	"archive/zip"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_sample.zip")
	source := map[string]string{
		"Hello": "world",
	}
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	writer, err := zipWriter.Create("hello.txt")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(writer, "%#v\n", source)
}

func main() {
	fmt.Println("port 8080")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
