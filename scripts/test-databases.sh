#!/bin/bash

echo "Testing Redis..."

# Test Redis
echo -n "Checking Redis... "
redis_response=$(docker compose exec redis redis-cli PING)
if [ "$redis_response" = "PONG" ]; then
    echo "✅"
    echo "Testing Redis write/read:"
    echo "Writing test value..."
    docker compose exec redis redis-cli SET test_key "Hello Redis"
    echo "Reading test value:"
    docker compose exec redis redis-cli GET test_key
else
    echo "❌ (Response: $redis_response)"
fi
