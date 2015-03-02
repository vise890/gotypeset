package main

import (
	"io"
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
	rawMmdIn, _, err := r.FormFile("inputMmd")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must select a file"))
		return
	}

	in, err := RegenerateFrontMatter(rawMmdIn)
	if err != nil {
		log.Fatal(err.Error())
	}
	out, err := mmd2pdf(in)
	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/pdf")
	_, err = io.Copy(w, out)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/typeset", typesetMarkdown)

	log.Print("Listening on :9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
