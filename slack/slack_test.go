package slack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Converter_Parse(t *testing.T) {

	converter := New()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "# Heading 1",
			expected: "*Heading 1*",
		},
		{
			input:    "## Heading 2",
			expected: "*Heading 2*",
		},
		{
			input:    "### Heading 3",
			expected: "*Heading 3*",
		},
		{
			input:    "#### Heading 4",
			expected: "*Heading 4*",
		},
		{
			input:    "##### Heading 5",
			expected: "*Heading 5*",
		},
		{
			input:    "###### Heading 6",
			expected: "*Heading 6*",
		},
		{
			input:    "above\n___\nbelow",
			expected: "above\n\n\n\nbelow",
		},
		{
			input:    "above\n\n---\nbelow",
			expected: "above\n\n\n\nbelow",
		},
		{
			input:    "above\n***\nbelow",
			expected: "above\n\n\n\nbelow",
		},
		{
			input:    "**This is bold text**",
			expected: "*This is bold text*",
		},
		{
			input:    "__This is bold text__",
			expected: "*This is bold text*",
		},
		{
			input:    "*This is italic text*",
			expected: "_This is italic text_",
		},
		{
			input:    "_This is italic text_",
			expected: "_This is italic text_",
		},
		{
			input:    "~~Strikethrough~~",
			expected: "~Strikethrough~",
		},
		{
			input:    "> blockquote",
			expected: "> blockquote",
		},
		{
			input:    "+ one\n+ two\n+ three",
			expected: "• one\n• two\n• three",
		},
		{
			input:    "- one\n- two\n- three",
			expected: "• one\n• two\n• three",
		},
		{
			input:    "* one\n* two\n* three",
			expected: "• one\n• two\n• three",
		},
		{
			input:    "1. one\n1. two\n1. three",
			expected: "1. one\n2. two\n3. three",
		},
		{
			input:    "1. one\n2. two\n3. three",
			expected: "1. one\n2. two\n3. three",
		},
		{
			// FIXME: issue with parser, should have start == 42
			input:    "42. one\n43. two\n44. three",
			expected: "1. one\n2. two\n3. three",
		},
		{
			// FIXME: issue with parser, should have start == 42
			input:    "42. one\n1. two\n1. three",
			expected: "1. one\n2. two\n3. three",
		},
		{
			input:    "`inline code`",
			expected: "`inline code`",
		},
		{
			input:    "```\ncode block\n```",
			expected: "```\ncode block\n```",
		},
		{
			input:    "\tindented code block\n\twith two lines",
			expected: "```\nindented code block\nwith two lines\n```",
		},
		{
			input:    "[evilmonkeyinc](https://github.com/evilmonkeyinc)",
			expected: "<https://github.com/evilmonkeyinc|evilmonkeyinc>",
		},
		{
			input: `
| Header 1 | Header 2 | Header 3 |
| --- | --- | --- |
| short value | longer value | really long value |
| qwerty | asdfgh | zxcvbn |
			`,
			expected: "*Header 1*   *Header 2*    *Header 3*\nshort value  longer value  really long value\nqwerty       asdfgh        zxcvbn",
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
			expected: "*Heading 1*\n\n*Heading 2*\n\n*Heading 3*\n\n*Heading 4*\n\n*Heading 5*\n\n*Heading 6*\n\n*This is bold text*\n\n\n*This is bold text*\n\n\n_This is italic text_\n\n\n_This is italic text_\n\n\n~Strikethrough~\n\n> blockquote\n• one\n• two\n• three\n\n1. one\n2. two\n3. three\n\n\n<https://github.com/evilmonkeyinc|evilmonkeyinc>\n\n*Header 1*   *Header 2*    *Header 3*\nshort value  longer value  really long value\nqwerty       asdfgh        zxcvbn",
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			actual, _ := converter.Parse([]byte(test.input))
			assert.Equal(t, test.expected, string(actual))
		})
	}
}
