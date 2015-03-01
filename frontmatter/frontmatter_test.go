package frontmatter

import (
	"testing"

	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func describeParsing(c gospec.Context) {
	c.Specify("parseFrontMatter", func() {
		c.Specify("parses a correctly formed frontMatter", func() {
			in := []byte("title: a tale of two gophers\nauthor: G. Gopherious")
			actual, err := parseFrontMatter(in)
			c.Assume(err, IsNil)

			expected := frontMatter{
				Title:  "a tale of two gophers",
				Author: "G. Gopherious",
			}

			c.Expect(actual, Equals, expected)
		})
		c.Specify("returns an error if title is missing", func() {
			in := []byte("author: G. Gopherious")
			actual, err := parseFrontMatter(in)

			expectedErr := ParseError{
				msg:  "you must specify a `title` (all lowercase) in your frontmatter",
				code: TitleRequiredErr,
			}

			c.Expect(actual, Equals, frontMatter{})
			c.Expect(err, Equals, expectedErr)
		})
		c.Specify("returns an error if author is missing", func() {
			in := []byte("title: a tale of two gophers")
			actual, err := parseFrontMatter(in)

			expectedErr := ParseError{
				msg:  "you must specify an `author` (all lowercase) in your frontmatter",
				code: AuthorRequiredError,
			}

			c.Expect(actual, Equals, frontMatter{})
			c.Expect(err, Equals, expectedErr)
		})
		c.Specify("returns an error the frontmatter is unparseable", func() {
			in := []byte("b;labber#blabbr")
			actual, _ := parseFrontMatter(in)

			c.Expect(actual, Equals, frontMatter{})
			// TODO: test that an error is actually returned!!
		})
	})
}

func TestFrontmatter(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(describeParsing)
	gospec.MainGoTest(r, t)
}
