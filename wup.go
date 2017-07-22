package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const defaultDest = "wup"

func copyTemp(prefix string, r io.Reader) (n int64,
	path string, err error) {
	tf, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return
	}
	defer tf.Close()

	path = tf.Name()

	n, err = io.Copy(tf, r)
	if err != nil {
		return
	}
	return
}

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

const uploadInfo = `to upload, map your revproxy /wup to localhost:%d,
then run: curl http://host/wup/dest --data-binary @src
and find your data in /tmp/dest
`
