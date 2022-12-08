package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	defaultDest = "wup"
	server      = "wup/CLR"
	overwriteHd = "X-Overwrite"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", server)

	_, dest := path.Split(r.URL.Path)
	if dest == "" {
		dest = defaultDest
	}
	destPath := filepath.Join(destDir, dest)
	if r.Header.Get(overwriteHd) != "yes" {
		if _, err := os.Stat(destPath); !os.IsNotExist(err) {
			http.Error(w,
				fmt.Sprint("FORBIDDEN cant overwrite existing file: ", destPath),
				403)
			return
		}
	}

	n, tempPath, err := uploadToTemp(dest, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprint("ERR cant upload to temp file: ", err), 500)
		return
	}
	if n == 0 {
		err = os.Remove(tempPath)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"WARN: wup cant remove zero-lenght temp file:", tempPath)
		}
		http.Error(w, "BAD zero-length input", 400)
		return
	}

	err = os.Rename(tempPath, destPath)
	if err != nil {
		http.Error(w,
			fmt.Sprint("ERR cant move uploaded file to dest path: ", err),
			500)
		return
	}

	w.WriteHeader(201)
	fmt.Fprintln(w, "CREATED", destPath)
}
