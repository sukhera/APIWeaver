# APIWeaver Scripts

Helper scripts for the APIWeaver project's design and development workflows.

## Scripts Overview

### `start-figma-mcp.sh`
Starts the Figma MCP server for integration with Claude CLI and Cursor.

**Usage:**
```bash
./scripts/start-figma-mcp.sh
```

**Features:**
- Automatically checks for and creates `.env` file if missing
- Validates Figma API token configuration
- Checks for port conflicts
- Provides colored output and clear error messages
- Handles graceful shutdown with Ctrl+C

**Prerequisites:**
- Node.js and npm installed
- Figma Personal Access Token (get from https://www.figma.com/settings)

### `claude-ux.sh`
Claude CLI wrapper for UX design workflows with design system integration.

**Usage:**
```bash
./scripts/claude-ux.sh "your prompt" [options]
```

**Options:**
- `--figma URL` - Include Figma design URL for analysis
- `--output PATH` - Save output to specified file
- `--personas` - Include user personas context
- `--tech` - Include tech stack context
- `--full-context` - Include all available context files
- `--no-design-system` - Skip design system context
- `--help` - Show help message

**Examples:**
```bash
# Basic component specification
./scripts/claude-ux.sh "Create button component specifications"

# Design review with Figma integration
./scripts/claude-ux.sh "Review this design for accessibility" --figma https://figma.com/file/abc123

# Comprehensive design with full context
./scripts/claude-ux.sh "Design modal system" --full-context --output docs/components/modal.md

# Design system audit
./scripts/claude-ux.sh "Audit design system consistency" --output audits/system-audit.md
```

## Setup Instructions

### 1. Initial Setup
```bash
# Make scripts executable (already done)
chmod +x scripts/*.sh

# Create .env file (will be created automatically by start-figma-mcp.sh)
cp .env.example .env  # if you have an example file
# OR run the start script and it will create a template
```

### 2. Configure Environment Variables
Edit `.env` file:
```bash
# Figma API Configuration
FIGMA_API_KEY=figd_your_actual_token_here
FIGMA_MCP_PORT=3845

# Claude CLI Configuration  
CLAUDE_API_KEY=your_claude_api_key_here

# Optional: Enable debug mode
DEBUG=false
FIGMA_MCP_DEBUG=false
```

### 3. Install Dependencies
```bash
# Install Claude CLI globally
npm install -g @anthropic-ai/claude-cli

# Install Figma MCP server globally
npm install -g figma-developer-mcp

# Configure Claude CLI
claude auth login
```

### 4. Verify Setup
```bash
# Test Figma MCP server
./scripts/start-figma-mcp.sh
# Should start without errors and show server info

# Test Claude UX script (in another terminal)
./scripts/claude-ux.sh "Test connection" --help
```

## Workflow Integration

### Daily Design Tasks
```bash
# Start Figma MCP server (run once per session)
./scripts/start-figma-mcp.sh &

# Create component specifications
./scripts/claude-ux.sh "Create data table component specs" --personas --output docs/components/

# Review existing designs
./scripts/claude-ux.sh "Review button variants for consistency" --figma FIGMA_URL

# Audit design system
./scripts/claude-ux.sh "Perform weekly design system audit" --full-context --output audits/
```

### Integration with Other Tools

#### Git Hooks
Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
# Run design system compliance check before commits
if git diff --cached --name-only | grep -q "docs/design-system.md\|src/components/"; then
  ./scripts/claude-ux.sh "Check design system compliance for modified files" \
    --output .git/design-compliance-check.md
fi
```

#### IDE Integration
Add to VS Code tasks (`.vscode/tasks.json`):
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Start Figma MCP",
      "type": "shell",
      "command": "./scripts/start-figma-mcp.sh",
      "group": "build",
      "presentation": {
        "echo": true,
        "reveal": "always",
        "panel": "new"
      }
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **"Claude CLI not found"**
   ```bash
   npm install -g @anthropic-ai/claude-cli
   claude auth login
   ```

2. **"Figma MCP server failed to start"**
   - Check your FIGMA_API_KEY in `.env`
   - Verify port 3845 is not in use: `lsof -i :3845`
   - Test token: `curl -H "X-Figma-Token: YOUR_TOKEN" https://api.figma.com/v1/me`

3. **"Context files not found"**
   - Ensure design system documentation exists: `docs/design-system.md`
   - Check agent configuration: `.claude/agent/ux-design-expert.md`
   - Verify context files: `.claude/context/`

4. **"Permission denied"**
   ```bash
   chmod +x scripts/*.sh
   ```

### Debug Mode
Enable debug output by setting environment variables:
```bash
export DEBUG=true
export FIGMA_MCP_DEBUG=true
./scripts/claude-ux.sh "your command"
```

### Log Files
Check logs for troubleshooting:
```bash
# Claude CLI logs
tail -f ~/.claude/logs/claude-cli.log

# Figma MCP server logs (if running with debug)
# Output will be in terminal where server was started
```

## Contributing

When adding new scripts:

1. Follow the existing pattern for error handling and colored output
2. Include comprehensive help messages
3. Add validation for prerequisites and inputs  
4. Update this README with usage instructions
5. Make scripts executable: `chmod +x scripts/new-script.sh`

## Related Documentation

- [Figma MCP + Claude CLI Guide](../docs/figma-mcp-claude-cli-guide.md) - Complete integration guide
- [Design System Documentation](../docs/design-system.md) - Design system specification
- [Agent Use Guide](../docs/agent-use-guide.md) - General agent usage patterns
