package convert

import (
	"github.com/jomei/notionapi"
)

// ToNotionBlocks converts markdown content to Notion blocks
func ToNotionBlocks(markdown string) ([]notionapi.Block, error) {
	// TODO: Implement markdown parsing using goldmark
	// 1. Parse markdown to AST
	// 2. Walk AST and convert nodes to Notion blocks
	// 3. Handle nested structures (lists, toggles)

	// Placeholder: return a simple paragraph
	blocks := []notionapi.Block{
		&notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Type:   notionapi.BlockTypeParagraph,
				Object: notionapi.ObjectTypeBlock,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: "TODO: Implement markdown parsing"}},
				},
			},
		},
	}

	return blocks, nil
}
