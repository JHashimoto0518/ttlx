# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-01

### Added

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
  - `--output` / `-o` flag for custom output path
  - `--dry-run` flag to print to stdout
- `ttlx validate` - Validate YAML configuration files
- `ttlx version` - Display version information

#### Configuration Features
- Profile-based SSH connection definitions
- Customizable timeout settings (global and per-profile)
- Command execution support at each hop
- Port configuration (default: 22)

#### Developer Experience
- Comprehensive test suite with 97.8% coverage
- Integration tests for end-to-end workflows
- golangci-lint configuration for code quality
- Detailed error messages with suggestions

#### Documentation
- Complete README.md with usage examples
- CONTRIBUTING.md with development guidelines
- YAML configuration schema documentation
- Example configuration files

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
- `internal/config`: 97.4%
- `internal/generator`: 98.0%
- Overall: 97.8%

### Breaking Changes
- Initial release, no breaking changes

### Known Limitations
- Global options `retry`, `log`, and `log_file` are defined but not yet implemented
- Manual Tera Term testing required for generated TTL scripts
- No support for interactive command execution within TTL scripts

### Security Considerations
- Passwords in YAML files should use environment variables or runtime prompts
- Direct password specification (`value:`) is supported but not recommended
- Generated TTL scripts may contain sensitive information

## [0.1.0] - Development

### Added
- Project initialization
- Basic Go module structure
- Initial development environment setup

[Unreleased]: https://github.com/JHashimoto0518/ttlx/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/JHashimoto0518/ttlx/releases/tag/v1.0.0
[0.1.0]: https://github.com/JHashimoto0518/ttlx/releases/tag/v0.1.0
