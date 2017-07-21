package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var fname = "/tmp/wupfile"

func handler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create(fname)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERR: %s\n", err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERR: %s\n", err)
		return
	}
	fmt.Fprintln(w, "OK")
}

func main() {
	port := os.Args[1]
	fmt.Fprintf(os.Stderr,
		"to upload, run: curl http://host:%s --data-binary @myfile\n",
		port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
