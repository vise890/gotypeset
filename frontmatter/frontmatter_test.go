package frontmatter

import (
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
