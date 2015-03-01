package frontmatter

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"text/template"

	"gopkg.in/yaml.v2"
)

const frontMatterBodySeparator string = "---\n"

// frontMatter represents the frontmatters that are expected as inputs
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
	articleTemplate, err := template.New("article.yaml").ParseFiles("frontmatter/templates/article.yaml")
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

// RegenerateFrontMatter re-creates the frontmatter of a MultiMarkDown document
// to make it ready for conversion to pdf with mmd2pdf
func RegenerateFrontMatter(mmdIn io.Reader) (fullMmd io.Reader, err error) {
	rawMmdIn, err := ioutil.ReadAll(mmdIn)
	inFrontmatter, body, err := splitOutFrontMatter(rawMmdIn)

	fullFrontMatter := toLaTeXFrontMatter(inFrontmatter)

	fullMmdB := bytes.Join([][]byte{fullFrontMatter, body}, []byte(frontMatterBodySeparator))

	fullMmd = bytes.NewReader(fullMmdB)
	return fullMmd, nil
}
