package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

const version = "1.1.0"

func main() {
	port := flag.String("p", "8080", "port to serve on")
	directory := flag.String("d", ".", "the directory of static files to serve")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	http.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "url query parameter is required", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "failed to fetch the resource", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
