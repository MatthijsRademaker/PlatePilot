# CLAUDE.md

## Pommel - Semantic Code Search

This sub-project (frontend) uses Pommel for semantic code search. Pommel indexes your codebase into semantic chunks (files, classes, methods) and enables natural language search.

**Supported languages** (full AST-aware chunking): Go, Java, C#, Python, JavaScript, TypeScript, JSX, TSX

### Code Search Priority

**IMPORTANT: Use `pm search` BEFORE using Grep/Glob for code exploration.**

When looking for:
- How something is implemented → `pm search "authentication flow"`
- Where a pattern is used → `pm search "error handling"`
- Related code/concepts → `pm search "database connection"`
- Code that does X → `pm search "validate user input"`

Only fall back to Grep/Glob when:
- Searching for an exact string literal (e.g., a specific error message)
- Looking for a specific identifier name you already know
- Pommel daemon is not running

### Quick Search Examples
```bash
# Search within this sub-project (default when running from here)
pm search "authentication logic"

# Search with JSON output
pm search "error handling" --json

# Search across entire monorepo
pm search "shared utilities" --all

# Search specific chunk levels
pm search "class definitions" --level class
```

### Available Commands
- `pm search <query>` - Search this sub-project (or use --all for everything)
- `pm status` - Check daemon status and index statistics
- `pm subprojects` - List all sub-projects
- `pm start` / `pm stop` - Control the background daemon

### Tips
- Searches default to this sub-project when you're in this directory
- Use `--all` to search across the entire monorepo
- Chunk levels: file (entire files), class (structs/interfaces/classes), method (functions/methods)
