#!/bin/bash

# Script to build frontend for APIWeaver
# TODO: Add Go embedding functionality when web server is implemented

set -e

# Color codes for better visual feedback
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🏗️  APIWeaver Frontend Build${NC}"
echo "================================"

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ package.json not found. Please run this script from the web directory${NC}"
    exit 1
fi

# Check prerequisites
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js not found. Please install Node.js 18+${NC}"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm not found. Please install npm${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Prerequisites met${NC}"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📦 Installing dependencies...${NC}"
    npm install
    echo -e "${GREEN}✅ Dependencies installed${NC}"
fi

# Build with progress
echo -e "${YELLOW}📦 Building React frontend...${NC}"
start_time=$(date +%s)
npm run build
end_time=$(date +%s)
build_time=$((end_time - start_time))

echo -e "${GREEN}✅ Build completed in ${build_time}s${NC}"

# Validate build output
if [ ! -f "dist/index.html" ]; then
    echo -e "${RED}❌ Build validation failed - index.html not found${NC}"
    exit 1
fi

# Calculate bundle size
if command -v du &> /dev/null; then
    BUNDLE_SIZE=$(du -sh dist/ | cut -f1)
    echo -e "${BLUE}📊 Bundle size: ${BUNDLE_SIZE}${NC}"
fi

# Display build summary
echo ""
echo -e "${BLUE}📊 Build Summary:${NC}"
echo "Build time: ${build_time}s"
if [ ! -z "$BUNDLE_SIZE" ]; then
    echo "Bundle size: ${BUNDLE_SIZE}"
fi

echo ""
echo -e "${PURPLE}📁 Build artifacts:${NC}"
ls -la dist/

echo ""
echo -e "${GREEN}✅ Frontend build completed successfully!${NC}"
echo ""
echo -e "${YELLOW}📝 TODO: Go embedding functionality will be added when web server is implemented${NC}"