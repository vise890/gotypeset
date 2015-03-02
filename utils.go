package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func withTempDir(f func() (io.Reader, error)) (io.Reader, error) {
	previousWd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current Working Directory", err)
	}
	defer os.Chdir(previousWd)

	tempDir, err := ioutil.TempDir("", "typesetForge")
	if err != nil {
		log.Fatal("Could not create temporary directory: "+tempDir, err)
	}
	defer os.RemoveAll(tempDir)

	err = os.Chdir(tempDir)
	if err != nil {
		log.Fatal("Could not cd to "+tempDir, err)
	}

	return f()
}
