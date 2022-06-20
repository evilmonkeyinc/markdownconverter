package http

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Converter_Format(t *testing.T) {
	actual := New().Format()
	assert.Equal(t, "http", actual)
}

func Test_Converter_Parse(t *testing.T) {

	converter := New()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "# Heading 1",
			expected: "<h1>Heading 1</h1>",
		},
		{
			input:    "## Heading 2",
			expected: "<h2>Heading 2</h2>",
		},
		{
			input:    "### Heading 3",
			expected: "<h3>Heading 3</h3>",
		},
		{
			input:    "#### Heading 4",
			expected: "<h4>Heading 4</h4>",
		},
		{
			input:    "##### Heading 5",
			expected: "<h5>Heading 5</h5>",
		},
		{
			input:    "###### Heading 6",
			expected: "<h6>Heading 6</h6>",
		},
		{
			input:    "above\n___\nbelow",
			expected: "<p>above</p>\n\n<hr>\n\n<p>below</p>",
		},
		{
			input:    "above\n\n---\nbelow",
			expected: "<p>above</p>\n\n<hr>\n\n<p>below</p>",
		},
		{
			input:    "above\n***\nbelow",
			expected: "<p>above</p>\n\n<hr>\n\n<p>below</p>",
		},
		{
			input:    "**This is bold text**",
			expected: "<p><strong>This is bold text</strong></p>",
		},
		{
			input:    "__This is bold text__",
			expected: "<p><strong>This is bold text</strong></p>",
		},
		{
			input:    "*This is italic text*",
			expected: "<p><em>This is italic text</em></p>",
		},
		{
			input:    "_This is italic text_",
			expected: "<p><em>This is italic text</em></p>",
		},
		{
			input:    "~~Strikethrough~~",
			expected: "<p><del>Strikethrough</del></p>",
		},
		{
			input:    "> blockquote",
			expected: "<blockquote>\n<p>blockquote</p>\n</blockquote>",
		},
		{
			input:    "+ one\n+ two\n+ three",
			expected: "<ul>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ul>",
		},
		{
			input:    "- one\n- two\n- three",
			expected: "<ul>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ul>",
		},
		{
			input:    "* one\n* two\n* three",
			expected: "<ul>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ul>",
		},
		{
			input:    "1. one\n1. two\n1. three",
			expected: "<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>",
		},
		{
			input:    "1. one\n2. two\n3. three",
			expected: "<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>",
		},
		{
			// FIXME: issue with parser, should have start == 42
			input:    "42. one\n43. two\n44. three",
			expected: "<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>",
		},
		{
			// FIXME: issue with parser, should have start == 42
			input:    "42. one\n1. two\n1. three",
			expected: "<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>",
		},
		{
			input:    "`inline code`",
			expected: "<p><code>inline code</code></p>",
		},
		{
			input:    "```\ncode block\n```",
			expected: "<pre><code>code block\n</code></pre>",
		},
		{
			input:    "\tindented code block\n\twith two lines",
			expected: "<pre><code>indented code block\nwith two lines\n</code></pre>",
		},
		{
			input:    "[evilmonkeyinc](https://github.com/evilmonkeyinc)",
			expected: "<p><a href=\"https://github.com/evilmonkeyinc\">evilmonkeyinc</a></p>",
		},
		{
			input: `
| Header 1 | Header 2 | Header 3 |
| --- | --- | --- |
| short value | longer value | really long value |
| qwerty | asdfgh | zxcvbn |
			`,
			expected: "<table>\n<thead>\n<tr>\n<th>Header 1</th>\n<th>Header 2</th>\n<th>Header 3</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>short value</td>\n<td>longer value</td>\n<td>really long value</td>\n</tr>\n\n<tr>\n<td>qwerty</td>\n<td>asdfgh</td>\n<td>zxcvbn</td>\n</tr>\n</tbody>\n</table>",
		},
		{
			input: `
# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

###### Heading 6

**This is bold text**

__This is bold text__

*This is italic text*

_This is italic text_

~~Strikethrough~~

> blockquote

* one
* two
* three

1. one
1. two
1. three

[evilmonkeyinc](https://github.com/evilmonkeyinc)

| Header 1 | Header 2 | Header 3 |
| --- | --- | --- |
| short value | longer value | really long value |
| qwerty | asdfgh | zxcvbn |
`,
			expected: "<h1>Heading 1</h1>\n\n<h2>Heading 2</h2>\n\n<h3>Heading 3</h3>\n\n<h4>Heading 4</h4>\n\n<h5>Heading 5</h5>\n\n<h6>Heading 6</h6>\n\n<p><strong>This is bold text</strong></p>\n\n<p><strong>This is bold text</strong></p>\n\n<p><em>This is italic text</em></p>\n\n<p><em>This is italic text</em></p>\n\n<p><del>Strikethrough</del></p>\n\n<blockquote>\n<p>blockquote</p>\n</blockquote>\n\n<ul>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ul>\n\n<ol>\n<li>one</li>\n<li>two</li>\n<li>three</li>\n</ol>\n\n<p><a href=\"https://github.com/evilmonkeyinc\">evilmonkeyinc</a></p>\n\n<table>\n<thead>\n<tr>\n<th>Header 1</th>\n<th>Header 2</th>\n<th>Header 3</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>short value</td>\n<td>longer value</td>\n<td>really long value</td>\n</tr>\n\n<tr>\n<td>qwerty</td>\n<td>asdfgh</td>\n<td>zxcvbn</td>\n</tr>\n</tbody>\n</table>",
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			actual, _ := converter.Parse([]byte(test.input))
			assert.Equal(t, test.expected, string(actual))
		})
	}
}
