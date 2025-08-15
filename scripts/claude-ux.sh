#!/bin/bash
# Claude UX helper script for APIWeaver project
# Usage: ./scripts/claude-ux.sh "your prompt" [additional options]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DESIGN_SYSTEM="docs/design-system.md"
DESIGN_SYSTEM_CONTEXT=".claude/context/design-system.md"
UX_AGENT=".claude/agent/ux-design-expert.md"
USER_PERSONAS=".claude/context/user-personas.md"
TECH_STACK=".claude/context/tech_stack.md"

# Function to check if file exists
check_file() {
    if [ ! -f "$1" ]; then
        echo -e "${RED}❌ File not found: $1${NC}"
        return 1
    fi
    return 0
}

# Function to check if Figma MCP server is running
check_figma_mcp() {
    if curl -s "http://127.0.0.1:3845/mcp" | grep -q "session ID"; then
        return 0
    else
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo -e "${BLUE}🎨 Claude UX Helper for APIWeaver${NC}"
    echo ""
    echo "Usage: $0 \"your prompt\" [options]"
    echo ""
    echo "Options:"
    echo "  --figma URL          Include Figma design URL"
    echo "  --output PATH        Output file path"
    echo "  --personas           Include user personas context"
    echo "  --tech               Include tech stack context"
    echo "  --full-context       Include all available context"
    echo "  --no-design-system   Skip design system context"
    echo "  --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 \"Create button component specs\""
    echo "  $0 \"Review this design\" --figma https://figma.com/file/abc123"
    echo "  $0 \"Design modal system\" --personas --output docs/components/modal.md"
    echo "  $0 \"Audit design system\" --full-context --output audits/system-audit.md"
    echo ""
    echo "Prerequisites:"
    echo "  - Claude CLI installed and configured"
    echo "  - Figma MCP server running (for --figma option)"
    echo "  - Design system documentation in place"
}

# Parse arguments
PROMPT=""
FIGMA_URL=""
OUTPUT_PATH=""
INCLUDE_PERSONAS=false
INCLUDE_TECH=false
INCLUDE_DESIGN_SYSTEM=true
FULL_CONTEXT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --figma)
            FIGMA_URL="$2"
            shift 2
            ;;
        --output)
            OUTPUT_PATH="$2"
            shift 2
            ;;
        --personas)
            INCLUDE_PERSONAS=true
            shift
            ;;
        --tech)
            INCLUDE_TECH=true
            shift
            ;;
        --full-context)
            FULL_CONTEXT=true
            INCLUDE_PERSONAS=true
            INCLUDE_TECH=true
            shift
            ;;
        --no-design-system)
            INCLUDE_DESIGN_SYSTEM=false
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            if [ -z "$PROMPT" ]; then
                PROMPT="$1"
            else
                echo -e "${RED}❌ Unknown option: $1${NC}"
                show_usage
                exit 1
            fi
            shift
            ;;
    esac
done

# Check if prompt is provided
if [ -z "$PROMPT" ]; then
    echo -e "${RED}❌ No prompt provided${NC}"
    show_usage
    exit 1
fi

echo -e "${GREEN}🎨 Running Claude UX command${NC}"
echo -e "${YELLOW}Prompt: $PROMPT${NC}"

# Check if Claude CLI is available
if ! command -v claude &> /dev/null; then
    echo -e "${RED}❌ Claude CLI not found${NC}"
    echo -e "${YELLOW}   Install with: npm install -g @anthropic-ai/claude-cli${NC}"
    exit 1
fi

# Build context arguments
CONTEXT_ARGS=""

# Always include design system context (unless explicitly disabled)
if [ "$INCLUDE_DESIGN_SYSTEM" = true ]; then
    if check_file "$DESIGN_SYSTEM"; then
        CONTEXT_ARGS="$CONTEXT_ARGS --context=$DESIGN_SYSTEM"
        echo -e "${GREEN}✅ Including design system context${NC}"
    elif check_file "$DESIGN_SYSTEM_CONTEXT"; then
        CONTEXT_ARGS="$CONTEXT_ARGS --context=$DESIGN_SYSTEM_CONTEXT"
        echo -e "${GREEN}✅ Including design system context (quick reference)${NC}"
    else
        echo -e "${YELLOW}⚠️  Design system context not found${NC}"
    fi
fi

# Include UX agent if available
if check_file "$UX_AGENT"; then
    CONTEXT_ARGS="$CONTEXT_ARGS --context=$UX_AGENT"
    echo -e "${GREEN}✅ Including UX design expert agent${NC}"
fi

# Include personas if requested
if [ "$INCLUDE_PERSONAS" = true ]; then
    if check_file "$USER_PERSONAS"; then
        CONTEXT_ARGS="$CONTEXT_ARGS --context=$USER_PERSONAS"
        echo -e "${GREEN}✅ Including user personas context${NC}"
    else
        echo -e "${YELLOW}⚠️  User personas context not found${NC}"
    fi
fi

# Include tech stack if requested
if [ "$INCLUDE_TECH" = true ]; then
    if check_file "$TECH_STACK"; then
        CONTEXT_ARGS="$CONTEXT_ARGS --context=$TECH_STACK"
        echo -e "${GREEN}✅ Including tech stack context${NC}"
    else
        echo -e "${YELLOW}⚠️  Tech stack context not found${NC}"
    fi
fi

# Check Figma MCP server if Figma URL is provided
if [ ! -z "$FIGMA_URL" ]; then
    echo -e "${YELLOW}🔗 Figma URL provided: $FIGMA_URL${NC}"
    if check_figma_mcp; then
        echo -e "${GREEN}✅ Figma MCP server is running${NC}"
        # Note: Figma URL will be included in the prompt, not as a separate argument
        # since Claude CLI doesn't have native Figma support yet
        PROMPT="$PROMPT

Figma Design Reference: $FIGMA_URL

Please analyze this Figma design as part of your response."
    else
        echo -e "${RED}❌ Figma MCP server not running${NC}"
        echo -e "${YELLOW}   Start it with: ./scripts/start-figma-mcp.sh${NC}"
        echo -e "${YELLOW}   Continuing without Figma integration...${NC}"
    fi
fi

# Build output arguments
OUTPUT_ARGS=""
if [ ! -z "$OUTPUT_PATH" ]; then
    # Create output directory if it doesn't exist
    OUTPUT_DIR=$(dirname "$OUTPUT_PATH")
    if [ ! -d "$OUTPUT_DIR" ]; then
        mkdir -p "$OUTPUT_DIR"
        echo -e "${GREEN}📁 Created output directory: $OUTPUT_DIR${NC}"
    fi
    OUTPUT_ARGS="--output=$OUTPUT_PATH"
    echo -e "${GREEN}📝 Output will be saved to: $OUTPUT_PATH${NC}"
fi

# Show context summary
echo -e "${BLUE}📋 Context Summary:${NC}"
echo "$CONTEXT_ARGS" | tr ' ' '\n' | grep -E "^--context=" | sed 's/--context=/  - /' || echo "  - No context files"

# Execute Claude command
echo -e "${GREEN}🚀 Executing Claude command...${NC}"
echo ""

# Build and execute the full command
FULL_COMMAND="claude \"$PROMPT\" $CONTEXT_ARGS $OUTPUT_ARGS"

# Show the command being executed (for debugging)
if [ "$DEBUG" = "true" ]; then
    echo -e "${YELLOW}Debug: $FULL_COMMAND${NC}"
fi

# Execute the command
eval $FULL_COMMAND

# Check if command was successful
if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Claude UX command completed successfully${NC}"
    if [ ! -z "$OUTPUT_PATH" ] && [ -f "$OUTPUT_PATH" ]; then
        echo -e "${GREEN}📄 Output saved to: $OUTPUT_PATH${NC}"
        echo -e "${YELLOW}   File size: $(stat -f%z "$OUTPUT_PATH" 2>/dev/null || stat -c%s "$OUTPUT_PATH" 2>/dev/null || echo "unknown") bytes${NC}"
    fi
else
    echo -e "${RED}❌ Claude UX command failed${NC}"
    exit 1
fi
