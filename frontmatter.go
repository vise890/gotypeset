package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"gopkg.in/yaml.v2"
)

const frontMatterSeparator string = "---\n"

// frontMatter represents the frontmatters that are expected as inputs
type frontMatter struct {
	Title  string
	Author string
}

var (
	// TitleRequired is returned if the input frontmatter does
	// not contain a `title`
	ErrTitleRequired = errors.New("you must specify a `title` (all lowercase) in your frontmatter")
	// AuthorRequired is returned if the input frontmatter does
	// not contain an `author`
	ErrAuthorRequired = errors.New("you must specify an `author` (all lowercase) in your frontmatter")
	// FrontMatterRequired is returned if the input does not
	// contain a frontmatter
	ErrFrontMatterRequired = errors.New("you must specify a frontmatter with a `title` and an `author` (all lowercase)")
)

func parseFrontMatter(in []byte) (frontMatter, error) {
	f := frontMatter{}
	err := yaml.Unmarshal(in, &f)
	if err != nil {
		return frontMatter{}, fmt.Errorf("Could not unmarshal frontatter: %s", err)
	}
	if f.Title == "" {
		return frontMatter{}, ErrTitleRequired
	}
	if f.Author == "" {
		return frontMatter{}, ErrAuthorRequired
	}
	return f, nil
}

func splitOutFrontMatter(mmdIn []byte) (f frontMatter, body []byte, err error) {
	parts := bytes.Split(mmdIn, []byte(frontMatterSeparator))
	if len(parts) == 1 {
		return frontMatter{}, []byte{}, ErrFrontMatterRequired
	}
	rawF := parts[0]
	body = bytes.Join(parts[1:], []byte(frontMatterSeparator))

	f, err = parseFrontMatter(rawF)
	if err != nil {
		return frontMatter{}, nil, err
	}

	return f, body, nil
}

func toLaTeXFrontMatter(inF frontMatter) (fullFrontMatter []byte, err error) {
	articleTemplate, err := template.ParseFiles("./frontmatter_templates/article.yaml")

	if err != nil {
		panic(err)
	}

	fullFrontMatterW := bytes.NewBuffer([]byte{})
	err = articleTemplate.Execute(fullFrontMatterW, inF)
	if err != nil {
		return nil, fmt.Errorf("Could not execute article template: %s", err)
	}

	fullFrontMatter, err = ioutil.ReadAll(fullFrontMatterW)
	if err != nil {
		return nil, err
	}

	return fullFrontMatter, nil
}

// RegenerateFrontMatter re-creates the frontmatter of a MultiMarkDown document
// to make it ready for conversion to pdf with mmd2pdf
func RegenerateFrontMatter(mmdIn io.Reader) (fullMmd io.Reader, err error) {
	rawMmdIn, err := ioutil.ReadAll(mmdIn)
	inFrontmatter, body, err := splitOutFrontMatter(rawMmdIn)
	if err != nil {
		return nil, fmt.Errorf("Could not split out frontmatter: %s", err)
	}

	fullFrontMatter, err := toLaTeXFrontMatter(inFrontmatter)
	if err != nil {
		return nil, fmt.Errorf("Could not convert frontmatter to a MMD/LaTeX one: %s", err)
	}

	fullMmdB := bytes.Join(
		[][]byte{
			fullFrontMatter,
			body},
		[]byte(frontMatterSeparator))

	fullMmd = bytes.NewReader(fullMmdB)
	return fullMmd, nil
}
