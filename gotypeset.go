package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v2"
)

const frontMatterBodySeparator string = "---\n"

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

type frontMatter struct {
	Title  string
	Author string
}

func parseFrontMatter(in []byte) (frontMatter, error) {
	f := frontMatter{}
	err := yaml.Unmarshal(in, &f)
	return f, err
}

func splitOutFrontMatter(mmdIn []byte) (f frontMatter, body []byte, err error) {
	parts := bytes.Split(mmdIn, []byte(frontMatterBodySeparator))
	rawF := parts[0]
	body = parts[1]
	f, err = parseFrontMatter(rawF)
	if err != nil {
		return frontMatter{}, nil, err
	}
	return f, body, nil
}

func toLaTeXFrontMatter(inF frontMatter) (fullFrontMatter []byte) {
	articleTemplate, err := template.New("article.yaml").ParseFiles("./frontmatters/article.yaml")
	if err != nil {
		panic(err)
	}

	fullFrontMatterW := bytes.NewBuffer([]byte{})
	err = articleTemplate.Execute(fullFrontMatterW, inF)
	if err != nil {
		log.Fatal("Could not generate a full template;", err)
	}

	fullFrontMatter, err = ioutil.ReadAll(fullFrontMatterW)
	if err != nil {
		log.Fatal("Could not read generated full template;", err)
	}

	return fullFrontMatter
}

func regenerateFrontMatter(mmdIn io.Reader) (fullMmd io.Reader, err error) {
	rawMmdIn, err := ioutil.ReadAll(mmdIn)
	inFrontmatter, body, err := splitOutFrontMatter(rawMmdIn)

	fullFrontMatter := toLaTeXFrontMatter(inFrontmatter)

	fullMmdB := bytes.Join([][]byte{fullFrontMatter, body}, []byte(frontMatterBodySeparator))

	fullMmd = bytes.NewReader(fullMmdB)
	return fullMmd, nil
}

func typesetMarkdown(w http.ResponseWriter, r *http.Request) {
	rawMmdIn, _, err := r.FormFile("inputMmd")
	if err != nil {
		log.Fatal(err.Error())
	}

	in, err := regenerateFrontMatter(rawMmdIn)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/typeset", typesetMarkdown)

	log.Print("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
