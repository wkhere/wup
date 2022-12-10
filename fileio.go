package main

import (
	"fmt"
	"io"
	"io/ioutil"
)

func uploadToTemp(prefix string, r io.Reader) (n int64, path string,
	err error) {
	tf, err := ioutil.TempFile(destDir, prefix)
	if err != nil {
		return
	}
	defer safeClose(tf, &err)

	path = tf.Name()
	limitR := &reader{io.LimitedReader{R: r, N: sizeLimit}}

	n, err = io.Copy(tf, limitR)
	if limitR.more() {
		io.Copy(ioutil.Discard, r)
		return n, path, errOverLimit
	}

	return n, path, nil
}

type reader struct {
	io.LimitedReader
}

func (r *reader) more() bool {
	var b [1]byte
	n, _ := r.R.Read(b[:])
	return n > 0
}

func safeClose(c io.Closer, errp *error) {
	cerr := c.Close()
	if *errp == nil {
		*errp = cerr
	}
}

var errOverLimit = fmt.Errorf("size over the limit of %d", sizeLimit)
