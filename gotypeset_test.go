package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func testGoTypeset(c gospec.Context) {
	c.Specify("toTex", func() {
		in := strings.NewReader("# Hello world")
		expecedOut := "\\part{Hello world}\n\\label{helloworld}\n"

		out, err := toTex(in)
		c.Assume(err, IsNil)

		actualOut, err := ioutil.ReadAll(out)
		c.Assume(err, IsNil)

		c.Expect(string(actualOut), Equals, expecedOut)
	})
}

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec(testGoTypeset)
	gospec.MainGoTest(r, t)
}
