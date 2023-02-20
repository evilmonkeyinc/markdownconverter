package slack

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// New returns a new instance of Converter
func New() *Converter {
	return &Converter{}
}

// Converter is the Slack markdwn Converter implementation
type Converter struct {
}

// Format returns a unique name for the converter
func (converter *Converter) Format() string {
	return "slack"
}

// Parse will parse the standard markdown and return the converted data
func (converter *Converter) Parse(markdwn []byte) ([]byte, error) {
	parser := parser.New()

	data := markdown.NormalizeNewlines(markdwn)
	node := parser.Parse(data)

	bytes := markdown.Render(node, &renderer{})
	return []byte(strings.TrimSpace(string(bytes))), nil
}

type renderer struct {
}

func (rend *renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	// fmt.Println(reflect.TypeOf(node), entering)
	if !entering {
		switch node.(type) {
		case *ast.Link, *ast.TableCell, *ast.TableBody, *ast.TableHeader:
			break
		default:
			fmt.Fprint(w, "\n")
			return ast.GoToNext
		}
		return ast.GoToNext
	}

	switch node.(type) {
	case *ast.BlockQuote:
		blockquote := node.(*ast.BlockQuote)
		for _, child := range blockquote.Children {
			childData := markdown.Render(child, rend)
			data := strings.TrimSpace(string(childData))
			fmt.Fprintf(w, "> %s", data)
		}
		return ast.SkipChildren
	case *ast.Code:
		code := node.(*ast.Code)
		codeBlock := strings.TrimSpace(string(code.Literal))
		fmt.Fprintf(w, "`%s`", codeBlock)
		return ast.GoToNext
	case *ast.CodeBlock:
		code := node.(*ast.CodeBlock)
		codeBlock := strings.TrimSpace(string(code.Literal))
		fmt.Fprintf(w, "```\n%s\n```", codeBlock)
		return ast.GoToNext
	case *ast.Del:
		strikethrough := node.(*ast.Del)
		for _, child := range strikethrough.Children {
			childData := markdown.Render(child, rend)
			clean := strings.TrimSpace(string(childData))
			fmt.Fprintf(w, "~%s~", clean)
		}
		return ast.SkipChildren
	case *ast.Document:
		return ast.GoToNext
	case *ast.Emph:
		emphasis := node.(*ast.Emph)
		for _, child := range emphasis.Children {
			childData := markdown.Render(child, rend)
			clean := strings.TrimSpace(string(childData))
			fmt.Fprintf(w, "_%s_", clean)
		}
		return ast.SkipChildren
	case *ast.Heading:
		heading := node.(*ast.Heading)
		for _, child := range heading.Children {
			childData := markdown.Render(child, rend)
			fmt.Fprintf(w, "\n*%s*", string(childData))
		}
		return ast.SkipChildren
	case *ast.HorizontalRule:
		fmt.Fprint(w, "\n\n")
		return ast.GoToNext
	case *ast.Link:
		link := node.(*ast.Link)
		destination := link.Destination
		title := ""
		for _, child := range link.Children {
			childData := string(markdown.Render(child, rend))
			title = strings.TrimSpace(childData)
		}

		// When you copy/paste maintaining it is right, but API is angle brackets
		// fmt.Fprintf(w, "[%s](%s)", title, string(destination))
		fmt.Fprintf(w, "<%s|%s>", string(destination), title)
		return ast.SkipChildren
	case *ast.List:
		list := node.(*ast.List)
		start := list.Start
		if start == 0 {
			start = 1
		}

		for idx, child := range list.Children {
			childData := markdown.Render(child, rend)
			clean := strings.TrimSpace(string(childData))
			prefix := "•"
			if list.ListFlags&ast.ListTypeOrdered != 0 {
				prefix = fmt.Sprintf("%d.", idx+start)
			}
			fmt.Fprintf(w, "%s %s\n", prefix, clean)
		}
		return ast.SkipChildren
	case *ast.ListItem:
		item := node.(*ast.ListItem)
		for _, child := range item.Children {
			childData := string(markdown.Render(child, rend))
			childData = strings.TrimSpace(childData)
			fmt.Fprint(w, childData)
		}
		return ast.SkipChildren
	case *ast.Paragraph:
		fmt.Fprint(w, "\n")
		return ast.GoToNext
	case *ast.Strong:
		strong := node.(*ast.Strong)
		for _, child := range strong.Children {
			childData := markdown.Render(child, rend)
			clean := strings.TrimSpace(string(childData))
			fmt.Fprintf(w, "*%s*", clean)
		}
		return ast.SkipChildren
	case *ast.Table:
		fmt.Fprint(w, "\n")
		tabWritter := tabwriter.NewWriter(w, 2, 2, 2, ' ', 0)
		table := node.(*ast.Table)

		var childData string = ""
		if len(table.Children) > 0 {
			headerNode := table.Children[0]
			childData += fmt.Sprintf("%s\n", string(markdown.Render(headerNode, rend)))
		}
		if len(table.Children) > 1 {
			bodyNode := table.Children[1].AsContainer()
			for _, child := range bodyNode.Children {
				childData += fmt.Sprintf("%s", string(markdown.Render(child, rend)))

			}
		}
		fmt.Fprintf(tabWritter, "%s\n", string(childData))
		return ast.SkipChildren
	case *ast.TableHeader:
		heading := node.(*ast.TableHeader)
		// First and only child should be a TableRow
		children := heading.Children[0].AsContainer().Children
		for idx, child := range children {
			if idx != 0 {
				fmt.Fprint(w, "\t")
			}
			childData := markdown.Render(child, rend)
			fmt.Fprintf(w, "*%s*", string(childData))
		}
		return ast.SkipChildren
	case *ast.TableRow:
		row := node.(*ast.TableRow)
		for idx, child := range row.Children {
			if idx != 0 {
				fmt.Fprint(w, "\t")
			}
			childData := markdown.Render(child, rend)
			fmt.Fprintf(w, "%s", string(childData))
		}
		return ast.SkipChildren
	case *ast.TableCell, *ast.TableBody:
		return ast.GoToNext
	case *ast.Text:
		text := node.(*ast.Text)
		fmt.Fprintf(w, "%s", string(text.Literal))
		return ast.GoToNext
	default:
		if leaf := node.AsLeaf(); leaf != nil {
			fmt.Fprintf(w, "%s", string(leaf.Literal))
		}

		if container := node.AsContainer(); container != nil {
			for _, child := range container.Children {
				childData := markdown.Render(child, rend)
				fmt.Fprintf(w, "%s", string(childData))
			}
			return ast.SkipChildren
		}
	}

	return ast.GoToNext
}

func (rend *renderer) RenderHeader(w io.Writer, ast ast.Node) {}

func (rend *renderer) RenderFooter(w io.Writer, ast ast.Node) {}
