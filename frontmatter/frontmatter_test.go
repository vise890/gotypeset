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
	})
}

func TestFrontmatter(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(describeParsing)
	gospec.MainGoTest(r, t)
}
