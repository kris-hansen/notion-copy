// Package notion provides a Notion API client wrapper
package notion

import (
	"context"
	"os"

	"github.com/jomei/notionapi"
)

// Client wraps the Notion API client
type Client struct {
	api *notionapi.Client
}

// New creates a new Notion client using NOTION_API_KEY from environment
func New() (*Client, error) {
	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return &Client{
		api: notionapi.NewClient(notionapi.Token(apiKey)),
	}, nil
}

// GetPage retrieves a page by ID
func (c *Client) GetPage(ctx context.Context, pageID string) (*notionapi.Page, error) {
	return c.api.Page.Get(ctx, notionapi.PageID(pageID))
}

// GetBlocks retrieves all blocks for a page or block
func (c *Client) GetBlocks(ctx context.Context, blockID string) ([]notionapi.Block, error) {
	var allBlocks []notionapi.Block
	var cursor notionapi.Cursor

	for {
		resp, err := c.api.Block.GetChildren(ctx, notionapi.BlockID(blockID), &notionapi.Pagination{
			StartCursor: cursor,
			PageSize:    100,
		})
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, resp.Results...)

		if !resp.HasMore {
			break
		}
		cursor = notionapi.Cursor(resp.NextCursor)
	}

	return allBlocks, nil
}

// CreatePage creates a new page under a parent
func (c *Client) CreatePage(ctx context.Context, parentID string, title string, blocks []notionapi.Block) (*notionapi.Page, error) {
	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(parentID),
		},
		Properties: notionapi.Properties{
			"title": notionapi.TitleProperty{
				Type: notionapi.PropertyTypeTitle,
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: title}},
				},
			},
		},
		Children: blocks,
	}

	return c.api.Page.Create(ctx, req)
}
