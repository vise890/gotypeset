package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func withTempDir(f func(string) (io.Reader, error)) (io.Reader, error) {

	tempDir, err := ioutil.TempDir("", "typesetForge")
	if err != nil {
		return nil, fmt.Errorf("Could not get/create a temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	return f(tempDir)
}
