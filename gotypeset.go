package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var mmd2pdfBin string

func init() {
	mmd, err := exec.LookPath("mmd2pdf")
	if err != nil {
		log.Fatal("multimarkdown (and mmd2pdf) need to be installed")
	}
	mmd2pdfBin = mmd
}

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

func mmd2pdf(mmdIn io.Reader) (pdfOut io.Reader, err error) {

	runMmd2pdf := func() (io.Reader, error) {
		inputFileName := "in.md"
		outputFileName := "in.pdf"

		inputF, err := os.Create(inputFileName)
		if err != nil {
			log.Fatal("Could not create file in "+inputFileName, err)
		}
		_, err = io.Copy(inputF, mmdIn)
		if err != nil {
			log.Fatal("Could not copy input markdown to disk", err)
		}

		yes := exec.Command("yes", "\n")
		mmd2pdf := exec.Command(mmd2pdfBin, inputFileName)
		mmd2pdf.Stdin, _ = yes.StdoutPipe()
		_ = mmd2pdf.Start()
		_ = yes.Start()
		_ = mmd2pdf.Wait()
		yes.Process.Kill()

		outputF, err := os.Open(outputFileName)
		if err != nil {
			log.Fatal("Could not open output file", err)
			return nil, err
		}

		return outputF, nil
	}

	return withTempDir(runMmd2pdf)
}

func typesetMarkdown(w http.ResponseWriter, r *http.Request) {
	in := r.Body
	out, err := mmd2pdf(in)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = io.Copy(w, out)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	http.HandleFunc("/typeset", typesetMarkdown)

	log.Print("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
