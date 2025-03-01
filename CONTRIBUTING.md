# Contributing to ChronoServe

We love your input! We want to make contributing to ChronoServe as easy and transparent as possible.

## Development Process

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Project Setup

```powershell
# Clone your fork
git clone https://github.com/YOUR-USERNAME/ChronoServe.git
cd ChronoServe

# Install dependencies
go mod download

# Run tests
make test

# Start development server
make dev
```

## Coding Standards

### Go Code

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `make lint` before submitting
- Write tests for new features
- Maintain test coverage above 80%
- Document all exported functions and types
- Use meaningful variable names
- Keep functions focused and small

### Documentation

- Update API documentation for any changes
- Maintain markdown formatting standards
- Include code examples where appropriate
- Keep security documentation current

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests liberally after the first line

Example:
```bash
Add password validation middleware

- Implement Argon2id password hashing
- Add constant-time comparison
- Update security documentation
- Add unit tests

Fixes #123
```

## Testing

```powershell
# Run all tests
make test

# Run specific tests
go test ./... -run TestYourFeature

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Pull Request Process

1. Update documentation for any changed functionality
2. Update the README.md with details of major changes
3. Add tests for new features
4. Ensure all tests pass and linting is clean
5. Request review from maintainers

### PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Security update

## Testing
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] Coverage maintained

## Documentation
- [ ] API docs updated
- [ ] README updated
- [ ] Security docs updated
```

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Accept constructive criticism
- Focus on what's best for the community
- Show empathy towards others

### Our Responsibilities

- Maintain code quality
- Review pull requests promptly
- Provide feedback constructively
- Keep discussions focused and productive

## Security

- Report security vulnerabilities privately
- Follow secure coding practices
- Update dependencies regularly
- Document security implications

## License

By contributing, you agree that your contributions will be licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

### License Header

Add this header to new source files:

```go
// Copyright (C) 2025 ToxicDev
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
```

## Questions?

Feel free to open an issue for:
- Usage questions
- Development questions
- Feature requests
- Bug reports