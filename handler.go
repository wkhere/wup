package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

const (
	defaultDest = "wup"
	server      = "wup/CLR"
)

func handler(w http.ResponseWriter, r *http.Request) {
	writeErr := func(err error) {
		http.Error(w, fmt.Sprint("ERR ", err), 500)
	}
	writeOK := func(msgs ...interface{}) {
		fmt.Fprintln(w, msgs...)
	}

	w.Header().Set("Server", server)

	_, dest := path.Split(r.URL.Path)
	if dest == "" {
		dest = defaultDest
	}

	n, tempPath, err := uploadToTemp(dest, r.Body)
	if err != nil {
		writeErr(fmt.Errorf("cannot upload to temp file: %s", err))
		return
	}
	if n == 0 {
		err = os.Remove(tempPath)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"WARN: wup could not remove temp file %s\n", tempPath)
		}
		writeOK("OK NOP")
		return
	}

	destPath := path.Join(destDir, dest)
	err = os.Rename(tempPath, destPath)
	if err != nil {
		writeErr(fmt.Errorf("cannot move uploaded file to dest path: %s", err))
		return
	}

	writeOK("OK", destPath)
}
