package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const defaultFile = "wupfile"

func handler(w http.ResponseWriter, r *http.Request) {
	writeErr := func(msg interface{}) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERR %s\n", msg)
	}

	file := strings.TrimLeft(r.URL.Path, "/")
	if file == "" {
		file = defaultFile
	}
	file = path.Join(os.TempDir(), file)
	f, err := os.Create(file)
	if err != nil {
		writeErr(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		writeErr(err)
		return
	}
	fmt.Fprintf(w, "OK %s\n", file)
}

func main() {
	port := os.Args[1]
	fmt.Fprintf(os.Stderr,
		"to upload, run: curl http://host:%s/dest --data-binary @src\n",
		port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
