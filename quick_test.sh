#!/bin/bash

echo "🚀 45AI Quick Test Script"
echo "========================"
echo ""

# Test backend health
echo "1️⃣ Testing Backend Health..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if [[ $HEALTH_RESPONSE == *"healthy"* ]]; then
    echo "✅ Backend is healthy!"
else
    echo "❌ Backend is not running. Start with: cd backend && go run cmd/api/main.go"
    exit 1
fi

# Test frontend
echo ""
echo "2️⃣ Testing Frontend..."
FRONTEND_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9000)
if [[ $FRONTEND_RESPONSE == "200" ]]; then
    echo "✅ Frontend is running!"
else
    echo "❌ Frontend is not running. Start with: cd frontend && npm run dev:h5"
    exit 1
fi

# Test authentication
echo ""
echo "3️⃣ Testing Authentication..."
AUTH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"code": "test_quicktest_123"}')

if [[ $AUTH_RESPONSE == *"token"* ]]; then
    echo "✅ Authentication is working!"
    echo "📊 User Created:"
    echo "$AUTH_RESPONSE" | jq '.user // .'
else
    echo "❌ Authentication failed:"
    echo "$AUTH_RESPONSE"
    exit 1
fi

echo ""
echo "🎉 All tests passed! Your 45AI app is ready for testing."
echo ""
echo "🔥 Next Steps:"
echo "   1. Open http://localhost:9000 in your browser"
echo "   2. Click 'Login with WeChat'"
echo "   3. Explore the gallery and profile pages"
echo "   4. Test the 'Coming Soon' feature modals"
echo ""
echo "📖 Full testing guide: See TESTING_GUIDE.md" 