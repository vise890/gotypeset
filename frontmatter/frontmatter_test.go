package frontmatter

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

func TestFrontMatterParsing(t *testing.T) {

	Convey("Given a valid frontmatter", t, func() {
		in := []byte("title: a tale of two gophers\nauthor: G. Gopherious")

		Convey("When it is parsed", func() {
			result, err := parseFrontMatter(in)

			Convey("The parse result should contain all the info", func() {
				So(result, ShouldResemble, frontMatter{
					Title:  "a tale of two gophers",
					Author: "G. Gopherious",
				})
			})

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a frontmatter without a Title", t, func() {
		in := []byte("author: G. Gopherious")

		Convey("When it is parsed", func() {
			result, err := parseFrontMatter(in)

			Convey("An appropriate error should be returned", func() {
				So(err, ShouldResemble, ParseError{
					msg:  "you must specify a `title` (all lowercase) in your frontmatter",
					code: TitleRequiredErr,
				})
			})

			Convey("The parse result should be empty", func() {
				So(result, ShouldResemble, frontMatter{})
			})
		})
	})

	Convey("Given a frontmatter without an Author", t, func() {
		in := []byte("title: a tale of two gophers")

		Convey("When it is parsed", func() {
			result, err := parseFrontMatter(in)

			Convey("An appropriate error should be returned", func() {
				So(err, ShouldResemble, ParseError{
					msg:  "you must specify an `author` (all lowercase) in your frontmatter",
					code: AuthorRequiredError,
				})
			})

			Convey("The result should be empty", func() {
				So(result, ShouldResemble, frontMatter{})
			})
		})
	})

	Convey("Given an invalid frontmatter", t, func() {
		in := []byte("b;labber#blabbr")

		Convey("When it is parsed", func() {
			result, err := parseFrontMatter(in)

			Convey("An appropriate error should be returned", func() {
				So(err, ShouldHaveSameTypeAs, &yaml.TypeError{})
			})

			Convey("The result should be empty", func() {
				So(result, ShouldResemble, frontMatter{})
			})
		})
	})
}

func TestFrontMatterSplitting(t *testing.T) {

	Convey("Given a markdown document with a valid frontmatter", t, func() {
		inFrontMatter := []byte("title: a tale of two gophers\nauthor: G. Gopherious")
		separator := []byte("\n---\n")
		inBody := []byte("This is ma' body")
		in := bytes.Join([][]byte{inFrontMatter, separator, inBody}, []byte(""))

		Convey("When the frontmatter is split out", func() {
			resultFrontMatter, resultBody, err := splitOutFrontMatter(in)

			Convey("The result should contain a parsed frontmatter", func() {
				So(resultFrontMatter, ShouldResemble, frontMatter{
					Title:  "a tale of two gophers",
					Author: "G. Gopherious",
				})
			})

			Convey("The result should contain the body", func() {
				So(resultBody, ShouldResemble, inBody)
			})

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given that there are multiple separators", t, func() {
		inFrontMatter := []byte("title: a tale of two gophers\nauthor: G. Gopherious")
		separator := []byte("\n---\n")
		inBody := bytes.Join(
			[][]byte{
				[]byte("My body"),
				separator,
				[]byte("Also part of my body. I shall be displeased if you chop it off")},
			[]byte(""))

		in := bytes.Join([][]byte{inFrontMatter, separator, inBody}, []byte(""))

		Convey("When the frontmatter is split out", func() {
			resultFrontMatter, resultBody, err := splitOutFrontMatter(in)

			Convey("The result should contain a parsed frontmatter", func() {
				So(resultFrontMatter, ShouldResemble, frontMatter{
					Title:  "a tale of two gophers",
					Author: "G. Gopherious",
				})
			})

			Convey("The result should contain the ENTIRE body (with separators)", func() {
				So(resultBody, ShouldResemble, inBody)
			})

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a markdown document without a frontmatter", t, func() {
		in := []byte("This is ma' body")

		Convey("When the  frontmatter is split out", func() {
			resultFrontMatter, resultBody, err := splitOutFrontMatter(in)

			Convey("An appropriate error should be returned", func() {
				So(err, ShouldResemble, ParseError{
					msg:  "you must specify a frontmatter with a `title` and an `author` (all lowercase)",
					code: FrontMatterRequiredError,
				})
			})

			Convey("The result should contain an empty frontmatter", func() {
				So(resultFrontMatter, ShouldResemble, frontMatter{})
			})

			Convey("The result should contain an empty body", func() {
				So(resultBody, ShouldResemble, []byte{})
			})
		})
	})
}
