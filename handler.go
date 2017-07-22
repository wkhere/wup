package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

const defaultDest = "wup"

func handler(w http.ResponseWriter, r *http.Request) {
	writeErr := func(msg interface{}) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERR %s\n", msg)
	}

	_, dest := path.Split(r.URL.Path)
	if dest == "" {
		dest = defaultDest
	}

	n, tempPath, err := copyTemp(dest, r.Body)
	if err != nil {
		writeErr(err)
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
		writeErr(err)
		return
	}

	fmt.Fprintf(w, "OK %s\n", destPath)
}
