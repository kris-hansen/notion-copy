# notion-copy

Copy Notion documents to the filesystem and back.

## Overview

`notion-copy` is a CLI tool for moving content between Notion and your local filesystem. No sync, no magic — just straightforward copy operations in either direction.

**Pull:** Notion pages → Markdown files  
**Push:** Markdown files → Notion pages

## Installation

```bash
go install github.com/kris-hansen/notion-copy@latest
```

Or build from source:

```bash
git clone https://github.com/kris-hansen/notion-copy
cd notion-copy
go build
```

## Configuration

### 1. Create a Notion Integration

1. Go to [notion.so/my-integrations](https://www.notion.so/my-integrations)
2. Click "New integration"
3. Give it a name (e.g., "notion-copy")
4. Copy the API key (starts with `secret_`)

### 2. Share Pages with Your Integration

In Notion, open the page you want to use and:
1. Click "..." menu → "Add connections"
2. Select your integration

### 3. Set Environment Variables

Create a `.env` file in `~/.config/notion-copy/` or your working directory:

```bash
# Required: Your Notion API key
NOTION_API_KEY=secret_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Optional: Default parent page for push command
# This lets you run `notion-copy push ./file.md` without specifying a page ID
NOTION_DEFAULT_PAGE_ID=abc123def456789
```

**Finding your page ID:** Open the page in Notion. The URL looks like:
```
https://notion.so/My-Page-Title-abc123def456789
                                 └── this is your page ID
```

## Usage

### Push (Filesystem → Notion)

Upload markdown files to Notion:

```bash
# Single file (uses NOTION_DEFAULT_PAGE_ID)
notion-copy push ./doc.md

# Single file to specific parent page
notion-copy push ./doc.md abc123def456

# Directory (all .md files)
notion-copy push ./docs/

# Recursive (subdirectories become nested pages)
notion-copy push ./docs/ -r

# Preview what would be created
notion-copy push ./docs/ --dry-run
```

**What gets converted:**
- Headings (H1, H2, H3) → Notion headings
- Paragraphs → Notion paragraphs
- Code blocks (with language) → Notion code blocks
- Block quotes → Notion quotes
- Horizontal rules → Notion dividers

**Page titles:** Extracted from the first `# Heading` in the file, or uses the filename.

### Pull (Notion → Filesystem)

Export Notion pages to markdown:

```bash
# Single page
notion-copy pull abc123def456 ./output/

# Include nested subpages
notion-copy pull abc123def456 ./output/ --recursive

# Include linked databases as CSV
notion-copy pull abc123def456 ./output/ --include-databases

# Flat structure (no subdirectories)
notion-copy pull abc123def456 ./output/ --flatten
```

## Block Mapping

| Notion Block | Markdown |
|--------------|----------|
| Heading 1 | `# Heading` |
| Heading 2 | `## Heading` |
| Heading 3 | `### Heading` |
| Paragraph | Plain text |
| Code | ` ```language ... ``` ` |
| Quote | `> text` |
| Divider | `---` |
| Bulleted list | `- item` |
| Numbered list | `1. item` |
| To-do | `- [ ] task` |

## Examples

### Push your notes folder

```bash
# Set your default page once
export NOTION_DEFAULT_PAGE_ID=abc123def456

# Push everything
notion-copy push ~/notes/ -r
```

### Export a Notion workspace

```bash
notion-copy pull <workspace-root-page-id> ./backup/ --recursive
```

### Preview before pushing

```bash
$ notion-copy push ./docs/ --dry-run
📄 [dry-run] ./docs/README.md → "Getting Started" (12 blocks)
📄 [dry-run] ./docs/api.md → "API Reference" (45 blocks)
📁 [dir] examples
📄 [dry-run] ./docs/examples/basic.md → "Basic Example" (8 blocks)
```

## License

Apache License 2.0
