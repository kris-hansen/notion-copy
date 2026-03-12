package convert

import (
	"strings"
	"testing"

	"github.com/jomei/notionapi"
)

func TestToMarkdown_Heading1(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.Heading1Block{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeHeading1},
			Heading1:   notionapi.Heading{RichText: []notionapi.RichText{{PlainText: "Title"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "# Title") {
		t.Errorf("expected '# Title', got '%s'", md)
	}
}

func TestToMarkdown_Heading2(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.Heading2Block{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeHeading2},
			Heading2:   notionapi.Heading{RichText: []notionapi.RichText{{PlainText: "Section"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "## Section") {
		t.Errorf("expected '## Section', got '%s'", md)
	}
}

func TestToMarkdown_Paragraph(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeParagraph},
			Paragraph:  notionapi.Paragraph{RichText: []notionapi.RichText{{PlainText: "Hello world"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "Hello world") {
		t.Errorf("expected 'Hello world', got '%s'", md)
	}
}

func TestToMarkdown_CodeBlock(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.CodeBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeCode},
			Code: notionapi.Code{
				Language: "go",
				RichText: []notionapi.RichText{{PlainText: "func main() {}"}},
			},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "```go") {
		t.Errorf("expected code block with go language, got '%s'", md)
	}
	if !strings.Contains(md, "func main() {}") {
		t.Errorf("expected code content, got '%s'", md)
	}
}

func TestToMarkdown_Quote(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.QuoteBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockQuote},
			Quote:      notionapi.Quote{RichText: []notionapi.RichText{{PlainText: "A wise quote"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "> A wise quote") {
		t.Errorf("expected '> A wise quote', got '%s'", md)
	}
}

func TestToMarkdown_Divider(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.DividerBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeDivider},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "---") {
		t.Errorf("expected '---', got '%s'", md)
	}
}

func TestToMarkdown_BulletedList(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.BulletedListItemBlock{
			BasicBlock:       notionapi.BasicBlock{Type: notionapi.BlockTypeBulletedListItem},
			BulletedListItem: notionapi.ListItem{RichText: []notionapi.RichText{{PlainText: "First item"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "- First item") {
		t.Errorf("expected '- First item', got '%s'", md)
	}
}

func TestToMarkdown_NumberedList(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.NumberedListItemBlock{
			BasicBlock:       notionapi.BasicBlock{Type: notionapi.BlockTypeNumberedListItem},
			NumberedListItem: notionapi.ListItem{RichText: []notionapi.RichText{{PlainText: "Step one"}}},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "1. Step one") {
		t.Errorf("expected '1. Step one', got '%s'", md)
	}
}

func TestToMarkdown_ToDo(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.ToDoBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeToDo},
			ToDo:       notionapi.ToDo{RichText: []notionapi.RichText{{PlainText: "Buy milk"}}, Checked: false},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "- [ ] Buy milk") {
		t.Errorf("expected '- [ ] Buy milk', got '%s'", md)
	}
}

func TestToMarkdown_ToDoChecked(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.ToDoBlock{
			BasicBlock: notionapi.BasicBlock{Type: notionapi.BlockTypeToDo},
			ToDo:       notionapi.ToDo{RichText: []notionapi.RichText{{PlainText: "Done task"}}, Checked: true},
		},
	}

	md, err := ToMarkdown(blocks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "- [x] Done task") {
		t.Errorf("expected '- [x] Done task', got '%s'", md)
	}
}

func TestRichTextToMarkdown_Bold(t *testing.T) {
	richText := []notionapi.RichText{
		{
			PlainText:   "bold text",
			Annotations: &notionapi.Annotations{Bold: true},
		},
	}

	result := richTextToMarkdown(richText)
	if result != "**bold text**" {
		t.Errorf("expected '**bold text**', got '%s'", result)
	}
}

func TestRichTextToMarkdown_Italic(t *testing.T) {
	richText := []notionapi.RichText{
		{
			PlainText:   "italic text",
			Annotations: &notionapi.Annotations{Italic: true},
		},
	}

	result := richTextToMarkdown(richText)
	if result != "*italic text*" {
		t.Errorf("expected '*italic text*', got '%s'", result)
	}
}

func TestRichTextToMarkdown_Code(t *testing.T) {
	richText := []notionapi.RichText{
		{
			PlainText:   "code",
			Annotations: &notionapi.Annotations{Code: true},
		},
	}

	result := richTextToMarkdown(richText)
	if result != "`code`" {
		t.Errorf("expected '`code`', got '%s'", result)
	}
}

func TestRichTextToMarkdown_Link(t *testing.T) {
	richText := []notionapi.RichText{
		{
			PlainText: "click here",
			Href:      "https://example.com",
		},
	}

	result := richTextToMarkdown(richText)
	if result != "[click here](https://example.com)" {
		t.Errorf("expected '[click here](https://example.com)', got '%s'", result)
	}
}
