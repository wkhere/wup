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
	writeErr := func(msg interface{}) {
		http.Error(w, fmt.Sprintf("ERR %s", msg), 500)
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
				"WARN: wup could not remove tempfile %s\n", tempPath)
		}
		fmt.Fprintf(w, "OK NOP\n")
		return
	}

	destPath := path.Join(os.TempDir(), dest)
	err = os.Rename(tempPath, destPath)
	if err != nil {
		writeErr(fmt.Errorf("cannot move uploaded file to dest path: %s", err))
		return
	}

	fmt.Fprintf(w, "OK %s\n", destPath)
}
