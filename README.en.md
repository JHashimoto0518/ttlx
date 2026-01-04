# ttlx - Tera Term Language eXtended

**Read this in other languages: [日本語](README.md) | English**

> [!WARNING]
> This project is under active development. APIs and functionality may change without notice.
> Not recommended for production use.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

ttlx is a tool that generates Tera Term macro (TTL) scripts from YAML configuration files.
By defining connection settings and commands as reusable profiles, it improves maintainability across multiple scripts.

## Tera Term Language Feature Support

| TTL Feature Category | Status | Description |
|---------------------|--------|-------------|
| **SSH Connection** | ✅ Supported | Multi-hop SSH connections (via bastion hosts) |
| **Authentication** | ✅ Supported | Password auth (password file/direct)<br>Public key authentication |
| **Command Execution** | ✅ Supported | Execute arbitrary commands after connection |
| **Error Handling** | ✅ Supported | Timeout handling, connection failure handling |
| **File Transfer** | 🔄 Not Yet | Planned for future release |
| **Dialog Display** | ⚠️ Partial | Password prompt and error messages only |
| **Variable Operations** | ⚠️ Partial | Password file reading, string concatenation |
| **Loops & Branching** | 🔄 Not Yet | Planned for future release |

## Features

- 📝 **YAML Configuration**: Define SSH routes and commands in a simple, readable YAML format
- 🔐 **Multiple Authentication Methods**: Support for password and public key authentication
- 🔗 **Multi-hop SSH**: Automate connections through bastion hosts and proxy servers
- ✅ **Validation**: Built-in configuration validation with helpful error messages
- 🎯 **Type-safe**: Leverages Go's type system for robust code generation
- 🧪 **Well-tested**: 97.8% test coverage

## Installation

### From Source

```bash
git clone https://github.com/JHashimoto0518/ttlx.git
cd ttlx
go build -o ttlx cmd/ttlx/main.go
```

### Using go install

```bash
go install github.com/JHashimoto0518/ttlx/cmd/ttlx@latest
```

## Quick Start

### 1. Create a YAML configuration file

```yaml
version: "1.0"

profiles:
  bastion:
    host: bastion.example.com
    user: user1
    prompt_marker: "$ "
    auth:
      type: password
      password_file: passwords.dat  # Read from password file

  target:
    host: 10.0.0.50
    user: user2
    prompt_marker: "$ "
    auth:
      type: password
      password_file: passwords.dat  # Read from password file
      password_prompt: "password:"  # Required for 2nd+ steps

routes:
  simple-connection:
    - profile: bastion
      commands:
        - echo "Connected to bastion"

    - profile: target
      commands:
        - ps aux
        - df -h
```

### 2. Generate TTL script

```bash
ttlx build config.yml
```

This creates `config.ttl` that you can run in Tera Term.

### 3. Validate configuration

```bash
ttlx validate config.yml
```

## Configuration

### Profile Settings

Each profile defines an SSH connection target:

```yaml
profiles:
  server_name:
    host: hostname_or_ip    # Required
    port: 22                # Optional, default: 22
    user: username          # Required
    prompt_marker: "$ "     # Required, prompt detection string
    auth:                   # Required
      type: password|keyfile
      # ... auth specific settings
```

### Authentication Types

#### Password Authentication

```yaml
auth:
  type: password
  # Choose one of the following:
  password_file: passwords.dat  # Read from password file (recommended)
                                # Default: "passwords.dat" if omitted
  value: "password"             # Direct password (for testing, not recommended)
  # Required for 2nd+ route steps:
  password_prompt: "password:"  # Password prompt string to wait for
```

**Using Password Files:**
- Uses Tera Term's `getpassword` command to read from password files
- Password name is automatically set to the profile name
- For creating password files, see [Tera Term Official Documentation](https://teratermproject.github.io/manual/5/en/macro/command/getpassword.html)

#### Public Key Authentication

```yaml
auth:
  type: keyfile
  path: ~/.ssh/id_rsa      # Path to private key file
```

### Route Configuration

Define the sequence of SSH connections:

```yaml
routes:
  multi-hop-connection:
    - profile: bastion       # First hop
      commands:              # Optional commands to run
        - su - root
        - cd /var/log

    - profile: target        # Second hop
      commands:
        - ps aux
```

### Global Options

```yaml
options:
  timeout: 30              # Connection timeout in seconds (default: 30)
  retry: 3                 # Number of retries (not implemented yet)
  log: true                # Enable logging (not implemented yet)
  log_file: /tmp/ttlx.log  # Log file path (not implemented yet)
  auto_disconnect: true    # Auto-disconnect after final step (default: false)
```

## CLI Commands

### build

Generate a TTL script from YAML configuration:

```bash
ttlx build <config.yml> [flags]

Flags:
  -o, --output string   Output directory path (default: current directory)
      --dry-run         Print to stdout instead of file

Example:
$ ttlx build config.yml
Generated TTL files:
  - config.ttl
  - config_simple-connection.ttl

$ ttlx build config.yml -o output/
Generated TTL files in output/:
  - config.ttl
  - config_simple-connection.ttl
```

### validate

Validate YAML configuration:

```bash
ttlx validate <config.yml>
```

### version

Print version information:

```bash
ttlx version
```

## Examples

See the [test/fixtures/valid](test/fixtures/valid) directory for more examples:

- [simple.yml](test/fixtures/valid/simple.yml) - Basic two-hop SSH connection
- [full.yml](test/fixtures/valid/full.yml) - Full-featured configuration with all options

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Building

```bash
go build -o ttlx cmd/ttlx/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter
golangci-lint run
```

### Project Structure

```
ttlx/
├── cmd/ttlx/          # CLI entry point
├── internal/
│   ├── cli/           # CLI commands
│   ├── config/        # Configuration handling
│   └── generator/     # TTL generation
├── test/
│   ├── fixtures/      # Test data
│   └── integration/   # Integration tests
└── docs/              # Documentation
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Tera Term Project](https://teratermproject.github.io/index-en.html) - The terminal emulator this tool generates scripts for
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-yaml](https://github.com/go-yaml/yaml) - YAML parser

## Support

- 🐛 [Report a bug](https://github.com/JHashimoto0518/ttlx/issues)
- 💡 [Request a feature](https://github.com/JHashimoto0518/ttlx/issues)
- 📖 [Documentation](docs/)
