package main

import (
	"io"
	"io/ioutil"
	"os"
)

func withTempDir(f func(string) (io.Reader, error)) (io.Reader, error) {

	tempDir, err := ioutil.TempDir("", "typesetForge")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	return f(tempDir)
}
