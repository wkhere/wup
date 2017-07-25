package main

import "os"

var destDir string

func init() {
	if destDir = os.Getenv("DESTDIR"); destDir == "" {
		destDir = os.TempDir()
	}
}
