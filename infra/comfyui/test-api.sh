#!/bin/bash

# Test script for ComfyUI API
# Usage: ./test-api.sh [API_URL]

API_URL=${1:-"http://localhost:8188"}

echo "Testing ComfyUI API at: $API_URL"
echo "=================================="

# Test 1: System Stats
echo "1. Testing system stats endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" "$API_URL/system_stats" 2>/dev/null)
http_code=${response##*HTTP_CODE:}
body=${response%HTTP_CODE:*}

if [ "$http_code" = "200" ]; then
    echo "✅ System stats endpoint working"
    echo "Response: $body" | jq . 2>/dev/null || echo "Response: $body"
else
    echo "❌ System stats endpoint failed (HTTP $http_code)"
    echo "Response: $body"
fi

echo

# Test 2: Queue endpoint
echo "2. Testing queue endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" "$API_URL/queue" 2>/dev/null)
http_code=${response##*HTTP_CODE:}
body=${response%HTTP_CODE:*}

if [ "$http_code" = "200" ]; then
    echo "✅ Queue endpoint working"
    echo "Queue length: $(echo "$body" | jq '.exec_info.queue_remaining // "N/A"' 2>/dev/null)"
else
    echo "❌ Queue endpoint failed (HTTP $http_code)"
    echo "Response: $body"
fi

echo

# Test 3: History endpoint
echo "3. Testing history endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" "$API_URL/history" 2>/dev/null)
http_code=${response##*HTTP_CODE:}

if [ "$http_code" = "200" ]; then
    echo "✅ History endpoint working"
else
    echo "❌ History endpoint failed (HTTP $http_code)"
fi

echo

# Test 4: Models endpoint
echo "4. Testing models endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" "$API_URL/object_info" 2>/dev/null)
http_code=${response##*HTTP_CODE:}
body=${response%HTTP_CODE:*}

if [ "$http_code" = "200" ]; then
    echo "✅ Models endpoint working"
    # Count available models
    checkpoint_count=$(echo "$body" | jq '.CheckpointLoaderSimple.input.required.ckpt_name[0] | length' 2>/dev/null || echo "N/A")
    echo "Available checkpoints: $checkpoint_count"
else
    echo "❌ Models endpoint failed (HTTP $http_code)"
fi

echo

# Test 5: Simple generation test (optional - requires workflow)
echo "5. Testing simple generation (if workflow exists)..."
if [ -f "test-workflow.json" ]; then
    echo "Found test workflow, running generation test..."
    response=$(curl -s -X POST -H "Content-Type: application/json" \
        -d @test-workflow.json \
        -w "HTTP_CODE:%{http_code}" \
        "$API_URL/prompt" 2>/dev/null)
    http_code=${response##*HTTP_CODE:}
    body=${response%HTTP_CODE:*}
    
    if [ "$http_code" = "200" ]; then
        echo "✅ Generation test started successfully"
        prompt_id=$(echo "$body" | jq -r '.prompt_id' 2>/dev/null)
        echo "Prompt ID: $prompt_id"
    else
        echo "❌ Generation test failed (HTTP $http_code)"
        echo "Response: $body"
    fi
else
    echo "ℹ️  No test workflow found, skipping generation test"
fi

echo
echo "=================================="
echo "API Test Complete"

# Summary
if curl -s "$API_URL/system_stats" >/dev/null 2>&1; then
    echo "✅ ComfyUI API is accessible and responding"
    echo "📍 API URL: $API_URL"
    echo "🔧 Web UI: $API_URL (if enabled)"
else
    echo "❌ ComfyUI API is not accessible"
    echo "🔍 Check if service is running and firewall allows port 8188"
fi 