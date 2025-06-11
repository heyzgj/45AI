#!/bin/bash

echo "🧪 45AI Authentication Testing Guide"
echo "=================================="
echo ""

echo "1. Testing Development Authentication:"
echo "   Frontend login will generate test codes like: test_1738097234_abc123"
echo ""

echo "2. Testing Backend Authentication API:"
echo "   curl -X POST http://localhost:8080/api/v1/auth/login \\"
echo "        -H 'Content-Type: application/json' \\"
echo "        -d '{\"code\": \"test_user1_12345\"}'"
echo ""

echo "3. Expected Response:"
echo "   {"
echo "     \"user\": {"
echo "       \"id\": 1,"
echo "       \"wechat_openid\": \"test_openid_user1_12345\","
echo "       \"nickname\": \"test_union_user1_12345\","
echo "       \"credits\": 50"
echo "     },"
echo "     \"token\": \"jwt_token_here\""
echo "   }"
echo ""

echo "🔥 Quick Test (press Enter to run):"
read -p ""
curl -X POST http://localhost:8080/api/v1/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"code": "test_user1_12345"}' | jq .

echo ""
echo "✅ If you see user data above, authentication is working!"
echo ""
echo "📱 To test the full frontend:"
echo "   1. Open http://localhost:9000"
echo "   2. Click 'Login with WeChat'"
echo "   3. It will automatically generate a test code and login"
echo ""
echo "🚀 For production WeChat testing:"
echo "   1. Use WeChat DevTools"
echo "   2. Import the frontend project as Mini Program"
echo "   3. Configure real WeChat App ID in backend .env" 