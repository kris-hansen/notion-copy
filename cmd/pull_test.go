package cmd

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Normal Title", "Normal Title"},
		{"Title/With/Slashes", "Title-With-Slashes"},
		{"Title: With Colon", "Title- With Colon"},
		{"Title?With*Special<Chars>", "TitleWithSpecialChars"},
		{"  Spaces  ", "Spaces"},
		{"...dots...", "dots"},
		{"", "untitled"},
		{"A very long title that exceeds one hundred characters and should be truncated to fit within the filename limit", "A very long title that exceeds one hundred characters and should be truncated to fit within the file"},
	}

	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractPageTitle(t *testing.T) {
	tests := []struct {
		name     string
		page     *notionapi.Page
		expected string
	}{
		{
			name: "title property",
			page: &notionapi.Page{
				Properties: notionapi.Properties{
					"title": &notionapi.TitleProperty{
						Title: []notionapi.RichText{{PlainText: "My Page"}},
					},
				},
			},
			expected: "My Page",
		},
		{
			name: "Name property",
			page: &notionapi.Page{
				Properties: notionapi.Properties{
					"Name": &notionapi.TitleProperty{
						Title: []notionapi.RichText{{PlainText: "Database Entry"}},
					},
				},
			},
			expected: "Database Entry",
		},
		{
			name: "no title",
			page: &notionapi.Page{
				Properties: notionapi.Properties{},
			},
			expected: "",
		},
		{
			name: "multi-part title",
			page: &notionapi.Page{
				Properties: notionapi.Properties{
					"title": &notionapi.TitleProperty{
						Title: []notionapi.RichText{
							{PlainText: "Part 1 "},
							{PlainText: "Part 2"},
						},
					},
				},
			},
			expected: "Part 1 Part 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPageTitle(tt.page)
			if result != tt.expected {
				t.Errorf("extractPageTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}
