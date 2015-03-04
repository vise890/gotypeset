package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vise890/gotypeset/frontmatter"
)

var mmd2pdfBin string

func init() {
	mmd, err := exec.LookPath("mmd2pdf")
	if err != nil {
		log.Fatal("multimarkdown (and mmd2pdf) need to be installed")
	}
	mmd2pdfBin = mmd
}

func mmd2pdf(mmdIn io.Reader) (pdfOut io.Reader, err error) {
	return withTempDir(func(workingDir string) (io.Reader, error) {
		inputFileName := filepath.Join(workingDir, "in.md")
		outputFileName := filepath.Join(workingDir, "in.pdf")

		inputF, err := os.Create(inputFileName)
		if err != nil {
			return nil, fmt.Errorf("Could not create input file: %s", err)
		}
		_, err = io.Copy(inputF, mmdIn)
		if err != nil {
			return nil, fmt.Errorf("Could not populate input file: %s", err)
		}

		yes := exec.Command("yes", "\n")
		mmd2pdf := exec.Command(mmd2pdfBin, inputFileName)
		mmd2pdf.Dir = workingDir

		mmd2pdf.Stdin, err = yes.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = mmd2pdf.Start()
		if err != nil {
			return nil, fmt.Errorf("Could not start mmd2pdf: %s", err)
		}
		err = yes.Start()
		if err != nil {
			return nil, err
		}
		err = mmd2pdf.Wait()
		if err != nil {
			return nil, err
		}
		yes.Process.Kill()

		outputF, err := os.Open(outputFileName)
		if err != nil {
			return nil, fmt.Errorf("Could not open output file: %s", err)
		}

		return outputF, nil
	})
}

func typesetMarkdown(w http.ResponseWriter, r *http.Request) {
	rawMmdIn, _, err := r.FormFile("inputMmd")
	if err != nil {
		log.Println(err)
		http.Error(w, "Must select a file.", http.StatusBadRequest)
		return
	}

	in, err := frontmatter.RegenerateFrontMatter(rawMmdIn)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not restructure document for typesetting.", http.StatusInternalServerError)
		return
	}
	out, err := mmd2pdf(in)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not typeset document.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	_, err = io.Copy(w, out)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not typeset document.", http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/typeset", typesetMarkdown)

	log.Print("Listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
