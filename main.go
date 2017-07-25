package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 9000, "")
	flag.Parse()

	fmt.Fprintf(os.Stderr, uploadInfo, *port, destDir)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const uploadInfo = `to upload, map your revproxy /wup to localhost:%d,
then run: curl http://host/wup/dest --data-binary @src
and find your data in %s/dest
`
