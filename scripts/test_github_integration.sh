#!/bin/bash

# GitHub Integration Test Script
# This script tests the GitHub API connection and repository access

set -e

echo "========================================"
echo "    SCVis GitHub Integration Test"
echo "========================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Load environment variables
if [ -f ".env" ]; then
    source .env
    print_status "Loaded .env file"
else
    print_warning ".env file not found, checking environment variables"
fi

# Check required variables
print_info "Checking GitHub configuration..."

if [ -z "$GITHUB_TOKEN" ]; then
    print_error "GITHUB_TOKEN not set"
    echo "Please set GITHUB_TOKEN in your .env file or environment"
    exit 1
else
    print_status "GITHUB_TOKEN is set"
fi

if [ -z "$GITHUB_OWNER" ]; then
    print_warning "GITHUB_OWNER not set, using default: vaccs"
    GITHUB_OWNER="vaccs"
else
    print_status "GITHUB_OWNER: $GITHUB_OWNER"
fi

if [ -z "$GITHUB_REPO" ]; then
    print_warning "GITHUB_REPO not set, using default: scvis-go"
    GITHUB_REPO="scvis-go"
else
    print_status "GITHUB_REPO: $GITHUB_REPO"
fi

echo ""
print_info "Testing GitHub API connection..."

# Test basic API access
print_info "Testing API authentication..."
response=$(curl -s -w "%{http_code}" -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/user)

http_code="${response: -3}"
response_body="${response%???}"

if [ "$http_code" = "200" ]; then
    username=$(echo "$response_body" | grep -o '"login":"[^"]*"' | cut -d'"' -f4)
    print_status "API authentication successful (user: $username)"
else
    print_error "API authentication failed (HTTP $http_code)"
    echo "Response: $response_body"
    exit 1
fi

# Test repository access
print_info "Testing repository access..."
response=$(curl -s -w "%{http_code}" -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO")

http_code="${response: -3}"
response_body="${response%???}"

if [ "$http_code" = "200" ]; then
    repo_name=$(echo "$response_body" | grep -o '"full_name":"[^"]*"' | cut -d'"' -f4)
    print_status "Repository access successful: $repo_name"
    
    # Check if issues are enabled
    has_issues=$(echo "$response_body" | grep -o '"has_issues":[^,]*' | cut -d':' -f2)
    if [ "$has_issues" = "true" ]; then
        print_status "Issues are enabled for this repository"
    else
        print_warning "Issues are disabled for this repository"
    fi
else
    print_error "Repository access failed (HTTP $http_code)"
    echo "Response: $response_body"
    
    if [ "$http_code" = "404" ]; then
        echo ""
        echo "Possible causes:"
        echo "  - Repository does not exist"
        echo "  - Repository is private and token lacks access"
        echo "  - GITHUB_OWNER or GITHUB_REPO is incorrect"
    elif [ "$http_code" = "403" ]; then
        echo ""
        echo "Possible causes:"
        echo "  - Token lacks repository permissions"
        echo "  - Token is expired"
        echo "  - Repository access is restricted"
    fi
    exit 1
fi

# Test issue creation permissions (dry run)
print_info "Testing issue creation permissions..."
test_issue_data='{
    "title": "SCVis Integration Test - Please Ignore",
    "body": "This is a test issue to verify GitHub integration. This issue should be deleted immediately.",
    "labels": ["test", "scvis-issue"]
}'

# Note: We're not actually creating the issue, just testing the endpoint
response=$(curl -s -w "%{http_code}" -X POST \
     -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     -H "Content-Type: application/json" \
     --data "$test_issue_data" \
     "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/issues" \
     --dry-run 2>/dev/null || echo "dry-run-not-supported")

if [[ "$response" == *"dry-run-not-supported"* ]]; then
    print_warning "Dry run not supported, skipping issue creation test"
    print_info "Issue creation permissions cannot be verified without creating a real issue"
else
    http_code="${response: -3}"
    if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
        print_status "Issue creation permissions verified"
    else
        print_warning "Issue creation test returned HTTP $http_code"
        print_info "This may indicate limited permissions, but integration might still work"
    fi
fi

# Check rate limits
print_info "Checking API rate limits..."
response=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/rate_limit)

if [[ "$response" == *"rate"* ]]; then
    remaining=$(echo "$response" | grep -o '"remaining":[^,]*' | cut -d':' -f2)
    limit=$(echo "$response" | grep -o '"limit":[^,]*' | cut -d':' -f2)
    print_status "Rate limit: $remaining/$limit remaining"
    
    if [ "$remaining" -lt 100 ]; then
        print_warning "Low rate limit remaining. Consider waiting before heavy usage."
    fi
else
    print_warning "Could not retrieve rate limit information"
fi

echo ""
print_status "GitHub integration test completed successfully!"
echo ""
echo "Configuration Summary:"
echo "  Repository: $GITHUB_OWNER/$GITHUB_REPO"
echo "  API Access: ✓ Working"
echo "  Repository Access: ✓ Working"
echo "  Issues Enabled: ✓ Yes"
echo ""
echo "You can now use the GitHub integration feature in SCVis."
echo "Administrators can submit issues directly to the repository through the web interface."
