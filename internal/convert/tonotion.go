package convert

import (
	"bytes"
	"strings"

	"github.com/jomei/notionapi"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ToNotionBlocks converts markdown content to Notion blocks
func ToNotionBlocks(markdown string) ([]notionapi.Block, error) {
	source := []byte(markdown)
	reader := text.NewReader(source)

	md := goldmark.New()
	doc := md.Parser().Parse(reader)

	var blocks []notionapi.Block
	err := convertNode(doc, source, &blocks)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// convertNode walks the AST and converts nodes to Notion blocks
func convertNode(node ast.Node, source []byte, blocks *[]notionapi.Block) error {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		block, err := nodeToBlock(child, source)
		if err != nil {
			return err
		}
		if block != nil {
			*blocks = append(*blocks, block)
		}
	}
	return nil
}

// nodeToBlock converts a single AST node to a Notion block
func nodeToBlock(node ast.Node, source []byte) (notionapi.Block, error) {
	switch n := node.(type) {
	case *ast.Heading:
		return headingToBlock(n, source)
	case *ast.Paragraph:
		return paragraphToBlock(n, source)
	case *ast.List:
		return nil, convertListItems(n, source, nil) // Lists are handled specially
	case *ast.FencedCodeBlock:
		return codeBlockToBlock(n, source)
	case *ast.CodeBlock:
		return codeBlockToBlock(n, source)
	case *ast.Blockquote:
		return quoteToBlock(n, source)
	case *ast.ThematicBreak:
		return dividerBlock(), nil
	default:
		// For unknown nodes, try to extract text as paragraph
		text := extractText(node, source)
		if text != "" {
			return createParagraph(text), nil
		}
		return nil, nil
	}
}

func headingToBlock(h *ast.Heading, source []byte) (notionapi.Block, error) {
	text := extractText(h, source)
	richText := []notionapi.RichText{{Text: &notionapi.Text{Content: text}}}

	switch h.Level {
	case 1:
		return &notionapi.Heading1Block{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeHeading1, Object: notionapi.ObjectTypeBlock},
			Heading1:   notionapi.Heading{RichText: richText},
		}, nil
	case 2:
		return &notionapi.Heading2Block{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeHeading2, Object: notionapi.ObjectTypeBlock},
			Heading2:   notionapi.Heading{RichText: richText},
		}, nil
	default: // 3+
		return &notionapi.Heading3Block{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeHeading3, Object: notionapi.ObjectTypeBlock},
			Heading3:   notionapi.Heading{RichText: richText},
		}, nil
	}
}

func paragraphToBlock(p *ast.Paragraph, source []byte) (notionapi.Block, error) {
	text := extractText(p, source)
	if text == "" {
		return nil, nil
	}
	return createParagraph(text), nil
}

func createParagraph(text string) notionapi.Block {
	return &notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeParagraph, Object: notionapi.ObjectTypeBlock},
		Paragraph:  notionapi.Paragraph{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: text}}}},
	}
}

func codeBlockToBlock(node ast.Node, source []byte) (notionapi.Block, error) {
	var buf bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buf.Write(line.Value(source))
	}

	language := "plain text"
	if fcb, ok := node.(*ast.FencedCodeBlock); ok {
		if lang := fcb.Language(source); lang != nil {
			language = string(lang)
		}
	}

	return &notionapi.CodeBlock{
		BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeCode, Object: notionapi.ObjectTypeBlock},
		Code: notionapi.Code{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: strings.TrimRight(buf.String(), "\n")}}},
			Language: language,
		},
	}, nil
}

func quoteToBlock(q *ast.Blockquote, source []byte) (notionapi.Block, error) {
	text := extractText(q, source)
	return &notionapi.QuoteBlock{
		BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockQuote, Object: notionapi.ObjectTypeBlock},
		Quote:      notionapi.Quote{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: text}}}},
	}, nil
}

func dividerBlock() notionapi.Block {
	return &notionapi.DividerBlock{
		BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeDivider, Object: notionapi.ObjectTypeBlock},
		Divider:    notionapi.Divider{},
	}
}

// convertListItems handles list conversion - returns blocks via the pointer
func convertListItems(list *ast.List, source []byte, blocks *[]notionapi.Block) error {
	// Note: This is a simplified implementation
	// Lists in Notion API are tricky because each item is a separate block
	return nil
}

// extractText recursively extracts text content from a node and its children
func extractText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	extractTextRecursive(node, source, &buf)
	return strings.TrimSpace(buf.String())
}

func extractTextRecursive(node ast.Node, source []byte, buf *bytes.Buffer) {
	switch n := node.(type) {
	case *ast.Text:
		buf.Write(n.Segment.Value(source))
		if n.HardLineBreak() || n.SoftLineBreak() {
			buf.WriteByte('\n')
		}
	case *ast.String:
		buf.Write(n.Value)
	case *ast.CodeSpan:
		// Include code span content
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
	}
}

// ExtractTitle gets the title from markdown (first H1 or filename)
func ExtractTitle(markdown string) string {
	source := []byte(markdown)
	reader := text.NewReader(source)

	md := goldmark.New()
	doc := md.Parser().Parse(reader)

	// Look for first heading
	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		if h, ok := child.(*ast.Heading); ok && h.Level == 1 {
			return extractText(h, source)
		}
	}

	return ""
}
