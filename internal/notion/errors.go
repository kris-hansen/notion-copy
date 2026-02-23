package notion

import "errors"

var (
	// ErrNoAPIKey is returned when NOTION_API_KEY is not set
	ErrNoAPIKey = errors.New("NOTION_API_KEY environment variable not set")
)
