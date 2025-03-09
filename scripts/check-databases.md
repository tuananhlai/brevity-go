# Testing Redis

## Using the Test Script
```bash
chmod +x scripts/test-databases.sh
./scripts/test-databases.sh
```

## Manual Testing

### Redis
1. Connect to Redis CLI:
```bash
docker compose exec redis redis-cli
```

2. Basic Redis commands:
```redis
PING                     # Should return PONG
SET mykey "Hello"       # Set a value
GET mykey              # Get the value
KEYS *                 # List all keys
```
