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

Set the following environment variables (or use a `.env` file):

```bash
NOTION_API_KEY=secret_xxx     # From notion.so/my-integrations
```

You can place your `.env` file in the working directory or `~/.config/notion-copy/.env`.

## Usage

### Pull (Notion → Filesystem)

Export a Notion page and its children to markdown:

```bash
notion-copy pull <page-id> ./output/
```

Options:
- `--recursive` — Include nested subpages
- `--include-databases` — Export linked databases as CSV
- `--flatten` — No subdirectories, flat file structure

### Push (Filesystem → Notion)

Import markdown files to Notion:

```bash
notion-copy push ./docs/ <parent-page-id>
```

Options:
- `--recursive` — Include subdirectories as nested pages
- `--dry-run` — Preview what would be created

## Block Mapping

| Notion Block | Markdown | Notes |
|--------------|----------|-------|
| Paragraph | Text | |
| Heading 1-3 | `#`, `##`, `###` | |
| Bulleted list | `- item` | |
| Numbered list | `1. item` | |
| To-do | `- [ ]` / `- [x]` | |
| Code | ` ``` ` | Language preserved |
| Quote | `>` | |
| Divider | `---` | |
| Image | `![alt](url)` | External URLs only |
| Toggle | `<details>` | HTML fallback |
| Callout | `> **Note:**` | Emoji prefix |
| Table | GFM table | |
| Bookmark | Link | URL only |

Unsupported blocks are exported as HTML comments with a warning.

## Examples

Back up your workspace:
```bash
notion-copy pull abc123def456 ./backup/$(date +%Y-%m-%d)/
```

Publish docs from a repo:
```bash
notion-copy push ./docs/ abc123def456 --recursive
```

## License

Apache License 2.0
