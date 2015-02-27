package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var multimarkdownBin string

func init() {
	mmd, err := exec.LookPath("multimarkdown")
	if err != nil {
		log.Fatal("multimarkdown needs to be installed")
	}
	multimarkdownBin = mmd
}

func mmd2tex(in io.Reader) (io.Reader, error) {
	cmd := exec.Command(multimarkdownBin, "--to=latex")
	cmd.Stdin = in
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}

func main() {
	inputPaths := os.Args[1:]
	if len(inputPaths) < 1 {
		log.Fatal("Must provide the path of at least one file to convert")
	}

	contents := make([][]byte, 0, 5)
	for _, p := range inputPaths {
		c, err := ioutil.ReadFile(p)
		if err != nil {
			log.Fatal("Could not read file at " + p)
		}
		contents = append(contents, c)
	}

	in := bytes.NewReader(bytes.Join(contents, []byte("\n\n")))
	out, err := mmd2tex(in)
	if err != nil {
		log.Fatal("Could not convert input to LaTeX")
	}
	output, err := ioutil.ReadAll(out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
