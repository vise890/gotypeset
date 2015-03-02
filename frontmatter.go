package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"text/template"

	"gopkg.in/yaml.v2"
)

const frontMatterSeparator string = "---\n"

// frontMatter represents the frontmatters that are expected as inputs
type frontMatter struct {
	Title  string
	Author string
}

type errorCode int

const (
	// TitleRequiredErr is returned if the input frontmatter does
	// not contain a `title`
	TitleRequiredErr errorCode = iota
	// AuthorRequiredError is returned if the input frontmatter does
	// not contain an `author`
	AuthorRequiredError
	// FrontMatterRequiredError is returned if the input does not
	// contain a frontmatter
	FrontMatterRequiredError
)

// ParseError is an error that can be returned when
// parsing frontmatters
type ParseError struct {
	msg  string
	code errorCode
}

func (e ParseError) Error() string {
	return e.msg
}

func parseFrontMatter(in []byte) (frontMatter, error) {
	f := frontMatter{}
	err := yaml.Unmarshal(in, &f)
	if err != nil {
		return frontMatter{}, err
	}
	if f.Title == "" {
		return frontMatter{}, ParseError{
			msg:  "you must specify a `title` (all lowercase) in your frontmatter",
			code: TitleRequiredErr,
		}
	}
	if f.Author == "" {
		return frontMatter{}, ParseError{
			msg:  "you must specify an `author` (all lowercase) in your frontmatter",
			code: AuthorRequiredError,
		}
	}
	return f, nil
}

func splitOutFrontMatter(mmdIn []byte) (f frontMatter, body []byte, err error) {
	parts := bytes.Split(mmdIn, []byte(frontMatterSeparator))
	if len(parts) == 1 {
		return frontMatter{}, []byte{}, ParseError{
			msg:  "you must specify a frontmatter with a `title` and an `author` (all lowercase)",
			code: FrontMatterRequiredError,
		}
	}
	rawF := parts[0]
	body = bytes.Join(parts[1:], []byte(frontMatterSeparator))

	f, err = parseFrontMatter(rawF)
	if err != nil {
		return frontMatter{}, nil, err
	}

	return f, body, nil
}

func toLaTeXFrontMatter(inF frontMatter) (fullFrontMatter []byte) {
	articleTemplate, err := template.New("article.yaml").ParseFiles("./frontmatter_templates/article.yaml")
	if err != nil {
		panic(err)
	}

	fullFrontMatterW := bytes.NewBuffer([]byte{})
	err = articleTemplate.Execute(fullFrontMatterW, inF)
	if err != nil {
		log.Fatal("Could not generate a full frontmatter;", err)
	}

	fullFrontMatter, err = ioutil.ReadAll(fullFrontMatterW)
	if err != nil {
		log.Fatal("Could not read generated full frontmatter;", err)
	}

	return fullFrontMatter
}

// RegenerateFrontMatter re-creates the frontmatter of a MultiMarkDown document
// to make it ready for conversion to pdf with mmd2pdf
func RegenerateFrontMatter(mmdIn io.Reader) (fullMmd io.Reader, err error) {
	rawMmdIn, err := ioutil.ReadAll(mmdIn)
	inFrontmatter, body, err := splitOutFrontMatter(rawMmdIn)
	if err != nil {
		return nil, err
	}

	fullFrontMatter := toLaTeXFrontMatter(inFrontmatter)

	fullMmdB := bytes.Join([][]byte{fullFrontMatter, body}, []byte(frontMatterSeparator))

	fullMmd = bytes.NewReader(fullMmdB)
	return fullMmd, nil
}
