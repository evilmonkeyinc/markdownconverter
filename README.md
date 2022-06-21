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

