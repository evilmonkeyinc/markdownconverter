[![codecov](https://codecov.io/gh/evilmonkeyinc/markdownconverter/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/markdownconverter)
[![Push Main](https://github.com/evilmonkeyinc/markdownconverter/actions/workflows/push_main.yml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/markdownconverter/actions/workflows/push_main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/evilmonkeyinc/markdownconverter.svg)](https://pkg.go.dev/github.com/evilmonkeyinc/markdownconverter)

# Markdown Converter

A simple command line tool that can perform markdown conversion to other formats

# Supported Formats

## Slack

A simple conversion between markdown and Slack markup (also known as `mrkdwn`).

Designed to be in the correct format for sending via the [Slack API](https://api.slack.com/methods/chat.postMessage) as `text` with `mrkdwn` set to true.

Slack `mrkdown` does not support all the features of markdown, as such some thing are not persisted perfectly such as different header levels or tables but this conversion should be enough for basic use cases such as posting a change-log or simple readme to a Slack message.

## HTML

A conversion between markdown and HTML, using the standard [gomarkdown/markdown](https://github.com/gomarkdown/markdown) `ToHTML` function with default options.

# Usage

## Command Line


```
Usage:

  markdownconverter [format] [input] [output]

Example:

  markdownconverter slack "[evilmonkeyinc](https://github.com/evilmonkeyinc)"
  > <https://github.com/evilmonkeyinc|evilmonkeyinc>

Options:

  -f, --format string   The output format
  -i, --input string    The input source file
  -o, --output string   The output destination file. optional
```

Download the latest version for your OS/Arch from the [Releases](https://github.com/evilmonkeyinc/markdownconverter/releases) page.

You can execute the tool from the command line with the following commands.
1. help - outputs the usage for the tool. You can also use the `--help`, or `-h` flag
2. version - outputs the version of the tool.
3. [format] [input] [output] - formats the input and returns it to the defined output file. If the output is not defined, it will be outputted to standard-out.

The arguments `format`, `input`, and `output` can be defined using flags with the same name if you want to change the order of arguments or just prefer using flags.

## Golang Module

Import `github.com/evilmonkeyinc/markdownconverter` into your golang project.

You can create a new instance of the `slack` or `http` converter by importing the specific package, and calling the exported `New()` function, which exposes the `Parse()` function which will take your markdown input and return the converted output.

```golang
...
import "github.com/evilmonkeyinc/markdownconverter/slack"
...

func main(){
    var inputBytes []byte
    // get your input
    ...
    converter := slack.New()
    outputBytes, err := converter.Parse(inputBytes)
    // use your output
}
```
