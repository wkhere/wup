package main

import (
	"io"
	"io/ioutil"
	"os"
)

func uploadToTemp(prefix string, r io.Reader) (n int64,
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
