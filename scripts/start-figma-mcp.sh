#!/bin/bash
# Start Figma MCP server for APIWeaver project
# Usage: ./scripts/start-figma-mcp.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
FIGMA_MCP_PORT=3845
ENV_FILE=".env"

echo -e "${GREEN}🎨 Starting Figma MCP Server for APIWeaver${NC}"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}⚠️  No .env file found. Creating template...${NC}"
    cat > "$ENV_FILE" << EOF
# Figma API Configuration
FIGMA_API_KEY=figd_YOUR_TOKEN_HERE
FIGMA_MCP_PORT=3845

# Claude CLI Configuration  
CLAUDE_API_KEY=your_claude_api_key_here

# Optional: Enable debug mode
# DEBUG=true
# FIGMA_MCP_DEBUG=true
EOF
    echo -e "${RED}❌ Please edit .env file with your actual Figma API token${NC}"
    echo -e "${YELLOW}   Get your token from: https://www.figma.com/settings${NC}"
    exit 1
fi

# Load environment variables
source "$ENV_FILE"

# Check if FIGMA_API_KEY is set
if [ -z "$FIGMA_API_KEY" ] || [ "$FIGMA_API_KEY" = "figd_YOUR_TOKEN_HERE" ]; then
    echo -e "${RED}❌ FIGMA_API_KEY not set in .env file${NC}"
    echo -e "${YELLOW}   Please edit .env file with your actual Figma API token${NC}"
    echo -e "${YELLOW}   Get your token from: https://www.figma.com/settings${NC}"
    exit 1
fi

# Check if port is already in use
if lsof -Pi :$FIGMA_MCP_PORT -sTCP:LISTEN -t >/dev/null ; then
    echo -e "${YELLOW}⚠️  Port $FIGMA_MCP_PORT is already in use${NC}"
    echo -e "${YELLOW}   Checking if it's our Figma MCP server...${NC}"
    
    # Test if it's responding to MCP requests
    if curl -s "http://127.0.0.1:$FIGMA_MCP_PORT/mcp" | grep -q "session ID"; then
        echo -e "${GREEN}✅ Figma MCP server is already running on port $FIGMA_MCP_PORT${NC}"
        exit 0
    else
        echo -e "${RED}❌ Port $FIGMA_MCP_PORT is occupied by another service${NC}"
        echo -e "${YELLOW}   Please stop the service or use a different port${NC}"
        exit 1
    fi
fi

# Check if figma-developer-mcp is installed
if ! command -v npx &> /dev/null; then
    echo -e "${RED}❌ npx not found. Please install Node.js${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Checking figma-developer-mcp installation...${NC}"
if ! npm list -g figma-developer-mcp &> /dev/null; then
    echo -e "${YELLOW}⚠️  figma-developer-mcp not installed globally. Installing...${NC}"
    npm install -g figma-developer-mcp
fi

# Start the server
echo -e "${GREEN}🚀 Starting Figma MCP server on port $FIGMA_MCP_PORT...${NC}"
echo -e "${YELLOW}   API Key: ${FIGMA_API_KEY:0:8}...${FIGMA_API_KEY: -3}${NC}"
echo -e "${YELLOW}   Port: $FIGMA_MCP_PORT${NC}"
echo -e "${YELLOW}   Press Ctrl+C to stop${NC}"
echo ""

# Set debug mode if enabled
if [ "$DEBUG" = "true" ] || [ "$FIGMA_MCP_DEBUG" = "true" ]; then
    export DEBUG=figma-mcp:*
    echo -e "${YELLOW}🐛 Debug mode enabled${NC}"
fi

# Start the server with error handling
trap 'echo -e "\n${YELLOW}🛑 Stopping Figma MCP server...${NC}"; exit 0' INT TERM

npx figma-developer-mcp \
    --figma-api-key="$FIGMA_API_KEY" \
    --port="$FIGMA_MCP_PORT" \
    --json=false \
    || {
        echo -e "${RED}❌ Failed to start Figma MCP server${NC}"
        echo -e "${YELLOW}   Check your API token and network connection${NC}"
        exit 1
    }
