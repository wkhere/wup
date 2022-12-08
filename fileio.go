package main

import (
	"io"
	"io/ioutil"
)

func uploadToTemp(prefix string, r io.Reader) (n int64, path string,
	err error) {
	tf, err := ioutil.TempFile(destDir, prefix)
	if err != nil {
		return
	}

	path = tf.Name()

	n, err = io.Copy(tf, r)

	if err2 := tf.Close(); err == nil {
		err = err2
	}
	return
}
