package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	plotter := exec.Command("dot", "-Tpng")

	source, e := plotter.StdinPipe()
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	png, e := plotter.StdoutPipe()
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	if e := plotter.Start(); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	defer plotter.Wait()

	go func() {
		defer source.Close() // NOTE: closing the pipe to notify the end of input.
		s, e := url.QueryUnescape(r.URL.RawQuery)
		if e != nil {
			log.Printf("Failed writing")
		}
		if _, e := source.Write([]byte(s)); e != nil {
			log.Printf("Failed writing")
		}
	}()

	if _, e := io.Copy(w, png); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
