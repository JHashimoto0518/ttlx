# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Note
Version 1.0.0 will be the first stable release. Currently in beta (0.1.0-beta).

## [0.1.0-beta] - Unreleased

### Added

#### Multiple Routes Support (Latest)
- **Multiple routes support**: Define multiple connection routes in a single YAML file
  - `routes` (plural) with named routes replaces single `route` field
  - Each route generates a separate TTL file
  - Profile reuse across multiple routes
  - Route name validation (alphanumeric, hyphens, underscores only)

#### Auto-Disconnect Option
- `auto_disconnect` option in global options for connection control
  - `true`: Automatically disconnect all connections and close Tera Term on success
  - `false` (default): Keep connections alive after script execution

#### Password Prompt Field
- `password_prompt` field for password authentication in multi-hop scenarios
  - Required for 2nd+ route steps with password authentication
  - Enables proper password input detection in nested SSH sessions
  - Security validation (rejects single quotes to prevent TTL injection)

#### Core Features
- YAML-based configuration for SSH connection routes
- Multi-hop SSH connection support through bastion hosts
- TTL script generation from YAML configuration
- Configuration validation with detailed error messages

#### Authentication
- Password authentication with multiple input methods:
  - Runtime password prompt (`prompt: true`)
  - Environment variable (`env: VAR_NAME`)
  - Direct password specification (`value: "password"`)
- Public key authentication support (`keyfile` type)

#### CLI Commands
- `ttlx build` - Generate TTL scripts from YAML configuration
  - `--output` / `-o` flag for custom output directory
  - `--dry-run` flag to print to stdout
- `ttlx validate` - Validate YAML configuration files
- `ttlx version` - Display version information

#### Configuration Features
- Profile-based SSH connection definitions
- Customizable timeout settings (global and per-profile)
- Command execution support at each hop
- Port configuration (default: 22)

#### Developer Experience
- Comprehensive test suite with 95.8% coverage
- Integration tests for end-to-end workflows
- golangci-lint configuration for code quality
- Detailed error messages with suggestions

#### Documentation
- Complete README.md with usage examples (Japanese and English)
- CONTRIBUTING.md with development guidelines
- YAML configuration schema documentation
- Example configuration files and test fixtures

### Changed

- **BREAKING**: `route` (singular) field is no longer supported. Use `routes` (plural) with named routes instead
- **BREAKING**: `-o` / `--output` flag now specifies output **directory** instead of output file path
  - Default: current directory (`.`)
  - Generated files are named `<route-name>.ttl`
  - Multiple TTL files are generated when multiple routes are defined
- CLI output format updated to show list of generated TTL files

### Fixed

- Include password in connect command `/passwd` option instead of separate `sendln` (#6)
- Fix connect result check from `<> 0` to `<> 2` (2 = linked and connected) (#8)
- Remove colon prefix from goto labels (`goto LABEL`, not `goto :LABEL`) (#8)
- Remove invalid timeout argument from wait command (#8)
- Add port option to SSH command (`ssh user@host -p port`) (#8)

### Technical Details

#### Project Structure
```
ttlx/
├── cmd/ttlx/          # CLI entry point
├── internal/
│   ├── cli/           # CLI commands (Cobra)
│   ├── config/        # Configuration handling
│   └── generator/     # TTL generation
└── test/              # Test fixtures and integration tests
```

#### Dependencies
- `github.com/spf13/cobra` v1.10.2 - CLI framework
- `gopkg.in/yaml.v3` v3.0.1 - YAML parsing
- `github.com/stretchr/testify` v1.11.1 - Testing utilities

#### Test Coverage
- `internal/config`: 96.8%
- `internal/generator`: 95.1%
- Overall: 95.8%

### Known Limitations
- Global options `retry`, `log`, and `log_file` are defined but not yet implemented
- Manual Tera Term testing required for generated TTL scripts
- No support for interactive command execution within TTL scripts

### Security Considerations
- Passwords in YAML files should use environment variables or runtime prompts
- Direct password specification (`value:`) is supported but not recommended
- Generated TTL scripts may contain sensitive information
- Route name validation prevents path traversal attacks

[Unreleased]: https://github.com/JHashimoto0518/ttlx/compare/v0.1.0-beta...HEAD
[0.1.0-beta]: https://github.com/JHashimoto0518/ttlx/releases/tag/v0.1.0-beta
