package main

import (
	"flag"
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
	port := flag.Int("port", 9000, "")
	flag.Parse()

	fmt.Fprintf(os.Stderr, uploadInfo, *port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const uploadInfo = `to upload, map your rev.proxy to localhost:%d,
then run: curl http://host/dest --data-binary @src
and get your file from /tmp/dest
`
