// DO NOT EDIT ** This file was generated with the bake tool ** DO NOT EDIT //

package frontmatter

var Templates = map[string]string{
	"article.yaml": `latex input:        mmd-article-header
Title:              {{.Title}}
Author:             {{.Author}}
Base Header Level:  2
LaTeX Mode:         memoir
latex input:        mmd-article-begin-doc
latex footer:       mmd-memoir-footer
`,

}
