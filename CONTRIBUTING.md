# Contributing to ChronoServe

We love your input! We want to make contributing to ChronoServe as easy and transparent as possible.

## Development Process

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Project Setup

```bash
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

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests liberally after the first line

## Testing

```bash
# Run all tests
make test

# Run specific tests
go test ./... -run TestYourFeature

# Run with race detection
go test -race ./...
```

## Pull Request Process

1. Update documentation for any changed functionality
2. Update the README.md with details of major changes
3. Add tests for new features
4. Ensure all tests pass and linting is clean
5. Request review from maintainers

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

## License

By contributing, you agree that your contributions will be licensed under the MIT License.