# Dev Swarm Project

This project is configured to use the **Dev Swarm** - an autonomous AI agent system for software development.

## Overview

The swarm consists of specialized AI agents that work together to implement features, review code, and manage the backlog:

| Agent | Role |
|-------|------|
| **Worker** | Implements features and fixes bugs |
| **Reviewer** | Reviews PRs and merges approved code |
| **Product Owner** | Generates and prioritizes stories |
| **Queen** | Designs feature specifications |

## Project Structure

```
.devagent/                  # Swarm project configuration
├── .env                    # Project ID and settings
├── vision.md               # Long-term project vision (guides PO agent)
├── architecture.md         # Auto-generated architecture docs
├── designs/                # Mockups, wireframes, diagrams
└── README.md               # This file

.claude/                    # Claude Code configuration
├── agents/                 # Specialized subagents for workers
│   ├── frontend-dev.md     # Vue/TypeScript/Tailwind specialist
│   ├── backend-dev.md      # Go/DDD specialist
│   ├── frontend-tester.md  # Vitest test specialist
│   ├── backend-tester.md   # Ginkgo BDD test specialist
│   └── e2e-tester.md       # Playwright E2E specialist
├── skills/                 # Claude Code skills
│   └── swarm-board/        # Task board interaction skill
└── settings.local.json     # MCP server configuration
```

## How It Works

### 1. Task Flow

```
Backlog → InProgress → Review → Done
              ↓          ↓
          (Worker)   (Reviewer)
              ↓          ↓
         Creates PR   Merges PR
```

### 2. Subagent Delegation

Workers delegate specialized work to subagents defined in `.claude/agents/`:

- **Frontend work** → `frontend-dev` subagent (Vue, Tailwind, Vuetify patterns)
- **Backend work** → `backend-dev` subagent (Go, DDD patterns)
- **Tests** → `*-tester` subagents (Vitest, Ginkgo, Playwright)

### 3. MCP Servers

The project is configured with MCP (Model Context Protocol) servers for enhanced capabilities:

- **tailwindcss** - Tailwind CSS class lookups and documentation
- **playwright** - Browser automation for E2E testing

MCP servers are pre-installed in the worker Docker image and configured in `.claude/settings.local.json`.

## Configuration Files

### vision.md

Edit `.devagent/vision.md` to guide the Product Owner agent:

- Define MVP scope and priorities
- Set short/medium/long-term goals
- Specify what's out of scope
- Describe target users

### Subagents

Customize `.claude/agents/*.md` for your project:

- Add project-specific coding conventions
- Update testing commands
- Add domain terminology
- Configure MCP tools

### CLAUDE.md

The root `CLAUDE.md` contains project-wide instructions for all Claude agents. This is merged with subagent prompts.

## Commands

```bash
# Task Management
swarm board add "Feature title"     # Create task (with Queen assistance)
swarm board add "Task" -d "desc"    # Create task directly
swarm board list                    # List all tasks
swarm board list --status Backlog   # Filter by status
swarm board show <id>               # Show task details

# Swarm Control
swarm up                            # Start 3 workers + 1 reviewer
swarm up 5                          # Start 5 workers + 1 reviewer
swarm up --enable-po                # Enable Product Owner agent
swarm down                          # Stop all containers
swarm status                        # Show swarm status

# Product Owner
swarm po start                      # Start PO daemon
swarm po status                     # Check PO status
swarm po stop                       # Stop PO daemon
```

## Dashboard

Access the task board dashboard at:

```
http://localhost:18080/dashboard
```

The dashboard shows:
- Task board with all statuses
- Active agents and their current tasks
- Real-time event stream
- Statistics and progress

## Manager API

The manager exposes a REST API at `http://localhost:18080/api/v1/`:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/tasks` | GET | List tasks |
| `/tasks` | POST | Create task |
| `/tasks/{id}` | GET | Get task details |
| `/stats` | GET | Dashboard statistics |
| `/agents` | GET | List active agents |
| `/events` | GET | Recent events |

## Customization

### Adding New Subagents

1. Create a new `.md` file in `.claude/agents/`
2. Use YAML frontmatter for configuration:
   ```yaml
   ---
   name: my-agent
   description: When to use this agent
   tools: Read, Edit, Write, Bash, Glob, Grep
   ---
   ```
3. Add your system prompt in the body
4. Commit the file

### Adding MCP Servers

1. Edit `.claude/settings.local.json`
2. Add new server configuration:
   ```json
   {
     "mcpServers": {
       "my-server": {
         "command": "bunx",
         "args": ["my-mcp-server"]
       }
     }
   }
   ```
3. Commit the file

## Troubleshooting

### Workers not finding subagents

Ensure `.claude/` is committed to the repository:
```bash
git add .claude/
git commit -m "chore: add Claude subagents"
```

### MCP servers not working

1. Check if MCP is configured: `cat .claude/settings.local.json`
2. Verify bun/bunx is available in the worker container
3. Check worker logs for MCP initialization errors

### Authentication issues

Workers share Claude authentication via Docker volume. Re-authenticate:
```bash
swarm up 1
ssh -p 2230 devuser@localhost  # Password: devpass
claude --login
exit
swarm down
```
