package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	defaultDest = "wup"
	server      = "wup/CLR"
	overwriteHd = "X-Overwrite"

	sizeLimit = 100 * 1024 * 1024
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", server)

	if n, _ := strconv.Atoi(r.Header.Get("Content-Length")); n > sizeLimit {
		// Error or missing header -> n is zero, which is ok for this check.
		// There will be another check when actually reading bytes.
		respError(w, 400, "BAD file size over the limit of ", sizeLimit)
		return
	}

	_, dest := path.Split(r.URL.Path)
	if dest == "" {
		dest = defaultDest
	}
	destPath := filepath.Join(destDir, dest)

	if _, err := os.Stat(destPath); !os.IsNotExist(err) {
		if r.Header.Get(overwriteHd) == "yes" {
			w.Header().Set(overwriteHd, "needed")
		} else {
			respError(w, 403, "FORBIDDEN cant overwrite file: ", destPath)
			return
		}
	}

	if r.Body == http.NoBody {
		respError(w, 400, "BAD zero-length input")
		return
	}

	_, tempPath, err := uploadToTemp(dest, r.Body)
	if err == errOverLimit {
		rerr := os.Remove(tempPath)
		if rerr != nil {
			respError(w, 500, "ERR cant remove interrupted temp file: ", rerr)
			return
		}
		respError(w, 400, "BAD upload interrupted: ", err)
		return
	}
	if err != nil {
		respError(w, 500, "ERR cant upload to temp file: ", err)
		return
	}

	err = os.Rename(tempPath, destPath)
	if err != nil {
		respError(w, 500, "ERR cant move uploaded file to dest path: ", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(201)
	fmt.Fprintln(w, "CREATED", destPath)
}

func respError(w http.ResponseWriter, code int, a ...interface{}) {
	http.Error(w, fmt.Sprint(a...), code)
}

func respErrorf(w http.ResponseWriter, code int, f string, a ...interface{}) {
	http.Error(w, fmt.Sprintf(f, a...), code)
}
