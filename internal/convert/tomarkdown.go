// Package convert handles conversion between Notion blocks and Markdown
package convert

import (
	"fmt"
	"strings"

	"github.com/jomei/notionapi"
)

// ToMarkdown converts Notion blocks to a markdown string
func ToMarkdown(blocks []notionapi.Block) (string, error) {
	var sb strings.Builder

	for _, block := range blocks {
		md, err := blockToMarkdown(block)
		if err != nil {
			// Write as comment for unsupported blocks
			sb.WriteString(fmt.Sprintf("<!-- Unsupported block type: %s -->\n", block.GetType()))
			continue
		}
		sb.WriteString(md)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func blockToMarkdown(block notionapi.Block) (string, error) {
	switch b := block.(type) {
	case *notionapi.ParagraphBlock:
		return richTextToMarkdown(b.Paragraph.RichText), nil

	case *notionapi.Heading1Block:
		return "# " + richTextToMarkdown(b.Heading1.RichText), nil

	case *notionapi.Heading2Block:
		return "## " + richTextToMarkdown(b.Heading2.RichText), nil

	case *notionapi.Heading3Block:
		return "### " + richTextToMarkdown(b.Heading3.RichText), nil

	case *notionapi.BulletedListItemBlock:
		return "- " + richTextToMarkdown(b.BulletedListItem.RichText), nil

	case *notionapi.NumberedListItemBlock:
		return "1. " + richTextToMarkdown(b.NumberedListItem.RichText), nil

	case *notionapi.ToDoBlock:
		checkbox := "[ ]"
		if b.ToDo.Checked {
			checkbox = "[x]"
		}
		return fmt.Sprintf("- %s %s", checkbox, richTextToMarkdown(b.ToDo.RichText)), nil

	case *notionapi.CodeBlock:
		lang := string(b.Code.Language)
		code := richTextToMarkdown(b.Code.RichText)
		return fmt.Sprintf("```%s\n%s\n```", lang, code), nil

	case *notionapi.QuoteBlock:
		return "> " + richTextToMarkdown(b.Quote.RichText), nil

	case *notionapi.DividerBlock:
		return "---", nil

	case *notionapi.ImageBlock:
		url := ""
		if b.Image.External != nil {
			url = b.Image.External.URL
		} else if b.Image.File != nil {
			url = b.Image.File.URL
		}
		caption := richTextToMarkdown(b.Image.Caption)
		return fmt.Sprintf("![%s](%s)", caption, url), nil

	case *notionapi.CalloutBlock:
		emoji := ""
		if b.Callout.Icon != nil && b.Callout.Icon.Emoji != nil {
			emoji = string(*b.Callout.Icon.Emoji) + " "
		}
		return fmt.Sprintf("> **%sNote:** %s", emoji, richTextToMarkdown(b.Callout.RichText)), nil

	case *notionapi.ToggleBlock:
		summary := richTextToMarkdown(b.Toggle.RichText)
		return fmt.Sprintf("<details>\n<summary>%s</summary>\n\n</details>", summary), nil

	case *notionapi.BookmarkBlock:
		return b.Bookmark.URL, nil

	default:
		return "", fmt.Errorf("unsupported block type: %s", block.GetType())
	}
}

func richTextToMarkdown(richText []notionapi.RichText) string {
	var parts []string
	for _, rt := range richText {
		text := rt.PlainText

		// Apply formatting
		if rt.Annotations != nil {
			if rt.Annotations.Code {
				text = "`" + text + "`"
			}
			if rt.Annotations.Bold {
				text = "**" + text + "**"
			}
			if rt.Annotations.Italic {
				text = "*" + text + "*"
			}
			if rt.Annotations.Strikethrough {
				text = "~~" + text + "~~"
			}
		}

		// Handle links
		if rt.Href != "" {
			text = fmt.Sprintf("[%s](%s)", text, rt.Href)
		}

		parts = append(parts, text)
	}
	return strings.Join(parts, "")
}
