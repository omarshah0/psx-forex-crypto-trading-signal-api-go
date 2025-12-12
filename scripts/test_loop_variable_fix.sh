#!/bin/bash

# Test Script: Verify Loop Variable Pointer Bug Fix
# This script tests that multi-package subscriptions return correct package details

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}  Testing Loop Variable Pointer Bug Fix${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
TOKEN="${AUTH_TOKEN:-}"

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ ERROR: AUTH_TOKEN environment variable not set${NC}"
    echo "Usage: AUTH_TOKEN=your_jwt_token ./scripts/test_loop_variable_fix.sh"
    exit 1
fi

echo -e "${GREEN}✓${NC} API URL: $API_URL"
echo -e "${GREEN}✓${NC} Auth token provided"
echo ""

# Test 1: Subscribe to multiple packages with different asset classes
echo -e "${YELLOW}Test 1: Subscribe to 3 different packages${NC}"
echo "Packages: [1] Forex Short Term, [5] Crypto Short Term, [9] PSX Long Term"
echo ""

RESPONSE=$(curl -s -X POST "$API_URL/api/subscriptions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "package_ids": [1, 5, 9]
  }')

echo "Response:"
echo "$RESPONSE" | jq '.' || echo "$RESPONSE"
echo ""

# Extract package IDs from subscriptions
PKG_ID_1=$(echo "$RESPONSE" | jq -r '.data.subscriptions[0].package.id // empty')
PKG_ID_2=$(echo "$RESPONSE" | jq -r '.data.subscriptions[1].package.id // empty')
PKG_ID_3=$(echo "$RESPONSE" | jq -r '.data.subscriptions[2].package.id // empty')

PKG_NAME_1=$(echo "$RESPONSE" | jq -r '.data.subscriptions[0].package.name // empty')
PKG_NAME_2=$(echo "$RESPONSE" | jq -r '.data.subscriptions[1].package.name // empty')
PKG_NAME_3=$(echo "$RESPONSE" | jq -r '.data.subscriptions[2].package.name // empty')

# Verify results
echo -e "${YELLOW}Verification:${NC}"

PASS_COUNT=0
FAIL_COUNT=0

# Check subscription 1
if [ "$PKG_ID_1" = "1" ]; then
    echo -e "${GREEN}✓${NC} Subscription 1: Package ID = $PKG_ID_1 (Expected: 1) - $PKG_NAME_1"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}✗${NC} Subscription 1: Package ID = $PKG_ID_1 (Expected: 1) - FAILED!"
    echo "   If this is 9, the loop variable bug still exists!"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Check subscription 2
if [ "$PKG_ID_2" = "5" ]; then
    echo -e "${GREEN}✓${NC} Subscription 2: Package ID = $PKG_ID_2 (Expected: 5) - $PKG_NAME_2"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}✗${NC} Subscription 2: Package ID = $PKG_ID_2 (Expected: 5) - FAILED!"
    echo "   If this is 9, the loop variable bug still exists!"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

# Check subscription 3
if [ "$PKG_ID_3" = "9" ]; then
    echo -e "${GREEN}✓${NC} Subscription 3: Package ID = $PKG_ID_3 (Expected: 9) - $PKG_NAME_3"
    PASS_COUNT=$((PASS_COUNT + 1))
else
    echo -e "${RED}✗${NC} Subscription 3: Package ID = $PKG_ID_3 (Expected: 9) - FAILED!"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi

echo ""

# Test 2: Check active subscriptions endpoint
echo -e "${YELLOW}Test 2: Verify active subscriptions endpoint${NC}"
echo ""

ACTIVE_RESPONSE=$(curl -s "$API_URL/api/subscriptions/active" \
  -H "Authorization: Bearer $TOKEN")

echo "Response:"
echo "$ACTIVE_RESPONSE" | jq '.' || echo "$ACTIVE_RESPONSE"
echo ""

# Count subscriptions
SUB_COUNT=$(echo "$ACTIVE_RESPONSE" | jq '.data | length')
echo "Active subscriptions found: $SUB_COUNT"

# Verify each has correct package
echo -e "${YELLOW}Verification:${NC}"

for i in $(seq 0 $((SUB_COUNT - 1))); do
    SUB_PKG_ID=$(echo "$ACTIVE_RESPONSE" | jq -r ".data[$i].package_id")
    RESP_PKG_ID=$(echo "$ACTIVE_RESPONSE" | jq -r ".data[$i].package.id")
    PKG_NAME=$(echo "$ACTIVE_RESPONSE" | jq -r ".data[$i].package.name")
    
    if [ "$SUB_PKG_ID" = "$RESP_PKG_ID" ]; then
        echo -e "${GREEN}✓${NC} Subscription package_id=$SUB_PKG_ID matches package.id=$RESP_PKG_ID ($PKG_NAME)"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}✗${NC} Subscription package_id=$SUB_PKG_ID DOES NOT match package.id=$RESP_PKG_ID - FAILED!"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}  Test Results${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "Passed: ${GREEN}$PASS_COUNT${NC}"
echo -e "Failed: ${RED}$FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ ALL TESTS PASSED!${NC}"
    echo "The loop variable pointer bug has been fixed."
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED!${NC}"
    echo "There may still be issues with the loop variable fix."
    exit 1
fi

