# TTL Test Environment

**Read this in other languages: [日本語](README.md) | English**

Test environment for Tera Term macros (TTL) using Docker Compose.

## Prerequisites

- Docker Desktop for Windows (WSL2 backend)
- Or Docker Engine (Linux/macOS)

## Test Procedures

### 1. Start Test SSH Servers

```bash
# Run in WSL2 or Linux environment
docker compose -f docker-compose.test.yml up -d

# Wait a bit for servers to start (takes longer on first run)
sleep 10

# Verify startup
docker compose -f docker-compose.test.yml ps
```

### 2. Connection Test (Linux Side)

```bash
# Test connection to bastion
ssh -p 2222 testuser@localhost
# Password: testpass123

# Test connection to target
ssh -p 2223 testuser@localhost
# Password: testpass123
```

### 3. Generate TTL Script

```bash
# Generate test TTL (specify directory path with -o)
./ttlx build test/test-config.yml -o test/

# Verify generated content
cat test/test-connection.ttl
```

### 4. Execute in Tera Term (Windows Side)

1. **Launch Tera Term**
2. **Run Macro**: Menu → Control → Macro
3. **Select File**: `\\wsl$\Ubuntu\workspaces\ttlx\test\test-connection.ttl`
   - Adjust path according to your environment
4. **Verify Operation**:
   - Can connect to bastion
   - Can multi-hop connect to target
   - Commands execute correctly

### 5. Clean Up Test Environment

```bash
# Stop and remove containers
docker compose -f docker-compose.test.yml down

# Also remove images
docker compose -f docker-compose.test.yml down --rmi all
```

## Test Configuration

### Connection Information

| Server | Host | Port | User | Password |
|--------|------|------|------|----------|
| bastion | localhost | 2222 | testuser | testpass123 |
| target  | localhost | 2223 | testuser | testpass123 |

**Note**: These credentials are for test environment only. Never use in production.

### Test Items

- [ ] Single-hop SSH connection (bastion)
- [ ] Multi-hop SSH connection (bastion → target)
- [ ] Password authentication (direct specification)
- [ ] Command execution
- [ ] Connection error handling
- [ ] Timeout handling

## Troubleshooting

### Port Already in Use

```bash
# Check port usage
netstat -an | grep -E "2222|2223"

# Use different ports by editing docker-compose.test.yml
```

### Container Won't Start

```bash
# Check logs
docker compose -f docker-compose.test.yml logs

# Restart containers
docker compose -f docker-compose.test.yml restart
```

### SSH Connection Fails

```bash
# Check SSH server status inside container
docker exec ttlx-test-bastion ps aux | grep ssh

# Check detailed SSH logs
ssh -vvv -p 2222 testuser@localhost
```

## Automation Script (Optional)

Automate from test environment startup to TTL generation:

```bash
#!/bin/bash
# test-setup.sh

echo "🚀 Setting up test environment..."

# 1. Start containers
docker compose -f docker-compose.test.yml up -d

# 2. Wait for servers to start
echo "⏳ Waiting for SSH servers to start..."
sleep 15

# 3. Test connections
echo "🔍 Testing connections..."
ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2222 testuser@localhost "echo 'bastion OK'" 2>/dev/null
ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2223 testuser@localhost "echo 'target OK'" 2>/dev/null

# 4. Generate TTL
echo "📝 Generating TTL script..."
./ttlx build test/test-config.yml -o test/

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Open Windows Explorer"
echo "2. Path: \\\\wsl$\\Ubuntu\\workspaces\\ttlx\\test\\test-connection.ttl"
echo "3. Execute macro in Tera Term"
```

Usage:

```bash
chmod +x test-setup.sh
./test-setup.sh
```
