package convert

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestToNotionBlocks_Heading1(t *testing.T) {
	markdown := "# Hello World"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	h1, ok := blocks[0].(*notionapi.Heading1Block)
	if !ok {
		t.Fatalf("expected Heading1Block, got %T", blocks[0])
	}

	if len(h1.Heading1.RichText) == 0 {
		t.Fatal("expected rich text content")
	}

	if h1.Heading1.RichText[0].Text.Content != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", h1.Heading1.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_Heading2(t *testing.T) {
	markdown := "## Section Title"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	h2, ok := blocks[0].(*notionapi.Heading2Block)
	if !ok {
		t.Fatalf("expected Heading2Block, got %T", blocks[0])
	}

	if h2.Heading2.RichText[0].Text.Content != "Section Title" {
		t.Errorf("expected 'Section Title', got '%s'", h2.Heading2.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_Heading3(t *testing.T) {
	markdown := "### Subsection"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	h3, ok := blocks[0].(*notionapi.Heading3Block)
	if !ok {
		t.Fatalf("expected Heading3Block, got %T", blocks[0])
	}

	if h3.Heading3.RichText[0].Text.Content != "Subsection" {
		t.Errorf("expected 'Subsection', got '%s'", h3.Heading3.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_Paragraph(t *testing.T) {
	markdown := "This is a paragraph of text."
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	p, ok := blocks[0].(*notionapi.ParagraphBlock)
	if !ok {
		t.Fatalf("expected ParagraphBlock, got %T", blocks[0])
	}

	if p.Paragraph.RichText[0].Text.Content != "This is a paragraph of text." {
		t.Errorf("expected paragraph text, got '%s'", p.Paragraph.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_CodeBlock(t *testing.T) {
	markdown := "```go\nfunc main() {}\n```"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	code, ok := blocks[0].(*notionapi.CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}

	if code.Code.Language != "go" {
		t.Errorf("expected language 'go', got '%s'", code.Code.Language)
	}

	if code.Code.RichText[0].Text.Content != "func main() {}" {
		t.Errorf("expected code content, got '%s'", code.Code.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_Quote(t *testing.T) {
	markdown := "> This is a quote"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	quote, ok := blocks[0].(*notionapi.QuoteBlock)
	if !ok {
		t.Fatalf("expected QuoteBlock, got %T", blocks[0])
	}

	if quote.Quote.RichText[0].Text.Content != "This is a quote" {
		t.Errorf("expected quote text, got '%s'", quote.Quote.RichText[0].Text.Content)
	}
}

func TestToNotionBlocks_Divider(t *testing.T) {
	markdown := "---"
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	_, ok := blocks[0].(*notionapi.DividerBlock)
	if !ok {
		t.Fatalf("expected DividerBlock, got %T", blocks[0])
	}
}

func TestToNotionBlocks_MultipleBlocks(t *testing.T) {
	markdown := `# Title

This is a paragraph.

## Section

More text here.

---

> A quote
`
	blocks, err := ToNotionBlocks(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have: H1, paragraph, H2, paragraph, divider, quote
	if len(blocks) < 5 {
		t.Errorf("expected at least 5 blocks, got %d", len(blocks))
	}

	// Check first block is H1
	if _, ok := blocks[0].(*notionapi.Heading1Block); !ok {
		t.Errorf("expected first block to be Heading1Block, got %T", blocks[0])
	}
}

func TestExtractTitle_WithH1(t *testing.T) {
	markdown := `# My Document Title

Some content here.
`
	title := ExtractTitle(markdown)
	if title != "My Document Title" {
		t.Errorf("expected 'My Document Title', got '%s'", title)
	}
}

func TestExtractTitle_WithoutH1(t *testing.T) {
	markdown := `## Not a title

Just a section.
`
	title := ExtractTitle(markdown)
	if title != "" {
		t.Errorf("expected empty string, got '%s'", title)
	}
}

func TestExtractTitle_EmptyDocument(t *testing.T) {
	markdown := ""
	title := ExtractTitle(markdown)
	if title != "" {
		t.Errorf("expected empty string, got '%s'", title)
	}
}
