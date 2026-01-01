# Contributing to ttlx

Thank you for your interest in contributing to ttlx! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- golangci-lint (for code quality checks)

### Setting Up Development Environment

1. Fork the repository on GitHub

2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/ttlx.git
cd ttlx
```

3. Add the upstream repository:
```bash
git remote add upstream https://github.com/JHashimoto0518/ttlx.git
```

4. Install dependencies:
```bash
go mod download
```

5. Verify the setup:
```bash
go test ./...
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### 2. Make Your Changes

- Follow the coding standards (see below)
- Write or update tests for your changes
- Update documentation as needed

### 3. Run Tests and Linters

```bash
# Run all tests
go test ./...

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run linter
golangci-lint run

# Format code
go fmt ./...
```

### 4. Commit Your Changes

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
git commit -m "type(scope): subject"
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

Examples:
```bash
git commit -m "feat(generator): add support for custom timeout per profile"
git commit -m "fix(validator): handle empty route array correctly"
git commit -m "docs(readme): add installation instructions"
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub with:
- Clear title and description
- Reference to any related issues
- Screenshots or examples if applicable

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Write clear, self-documenting code
- Add comments for exported functions and types

### Package Comments

All packages must have a package-level comment:

```go
// Package config handles YAML configuration file loading, validation, and data models.
package config
```

### Error Handling

- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf`:
```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

### Testing

- Write table-driven tests when possible
- Aim for high test coverage (>80%)
- Include both positive and negative test cases
- Use `testify` for assertions

Example:
```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  validConfig(),
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Project Structure

```
ttlx/
├── cmd/ttlx/          # CLI entry point
│   └── main.go
├── internal/          # Private application code
│   ├── cli/           # CLI commands (Cobra)
│   ├── config/        # Configuration handling
│   │   ├── model.go       # Data structures
│   │   ├── loader.go      # YAML loading
│   │   └── validator.go   # Validation logic
│   └── generator/     # TTL generation
│       ├── template.go    # TTL templates
│       └── generator.go   # Generation logic
├── test/
│   ├── fixtures/      # Test data (YAML files)
│   └── integration/   # Integration tests
└── docs/              # Documentation
```

## Adding New Features

### 1. Configuration Changes

If adding new configuration fields:

1. Update `internal/config/model.go`
2. Add validation in `internal/config/validator.go`
3. Update tests in `internal/config/*_test.go`
4. Add test fixtures in `test/fixtures/`

### 2. TTL Generation Changes

If modifying TTL output:

1. Update templates in `internal/generator/template.go`
2. Update generation logic in `internal/generator/generator.go`
3. Add tests in `internal/generator/generator_test.go`
4. Update expected output in test fixtures

### 3. CLI Changes

If adding new commands or flags:

1. Add command file in `internal/cli/`
2. Register in `internal/cli/root.go`
3. Update README.md documentation

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update CHANGELOG.md (see below)

## Pull Request Checklist

Before submitting a pull request, ensure:

- [ ] Code follows the style guidelines
- [ ] All tests pass (`go test ./...`)
- [ ] Test coverage is maintained or improved
- [ ] golangci-lint passes with no errors
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits
- [ ] CHANGELOG.md is updated (for significant changes)

## Release Process

(For maintainers)

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push a git tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```
4. GitHub Actions will handle the release build

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions in GitHub Discussions
- Contact maintainers via email (see README.md)

## License

By contributing to ttlx, you agree that your contributions will be licensed under the MIT License.
