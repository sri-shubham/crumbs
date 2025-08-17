# Contributing to Crumbs

Thank you for considering contributing to Crumbs! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by its Code of Conduct. Please be respectful and considerate of others.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Crumbs.

#### Before Submitting A Bug Report

- Check the [issues](https://github.com/sri-shubham/crumbs/issues) to see if the bug has already been reported
- Try to reproduce the issue with the latest version of the code
- Gather information about your environment (Go version, OS, etc.)

#### How Do I Submit A Good Bug Report?

Bugs are tracked as [GitHub issues](https://github.com/sri-shubham/crumbs/issues). Create an issue and provide the following information:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples to demonstrate the steps
- Describe the behavior you observed after following the steps
- Explain which behavior you expected to see instead and why
- Include environment details (Go version, OS, etc.)
- Include code samples or test cases that demonstrate the issue

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Crumbs.

#### Before Submitting An Enhancement Suggestion

- Check the [issues](https://github.com/sri-shubham/crumbs/issues) for similar enhancement suggestions
- Check if the enhancement has already been implemented in the latest version
- Consider whether your idea fits within the scope and aims of the project

#### How Do I Submit A Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/sri-shubham/crumbs/issues). Create an issue and provide the following information:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful to most users
- Provide examples of how the enhancement would be used
- Consider including code examples showing what would change

### Pull Requests

The process described here has several goals:

- Maintain the project's quality
- Fix problems that are important to users
- Enable a sustainable system for the project's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

1. **Fork the repository** and create your branch from `main`
2. **Make your changes**:
   - Follow the coding style used in the project
   - Add tests for any new functionality
   - Update documentation as needed
3. **Ensure the test suite passes**:
   - Run `go test ./...` to run all tests
   - Run `go test -bench=. -benchmem` to ensure performance hasn't been affected
4. **Create a pull request** with a clear title and description

#### Pull Request Requirements

- The code must be formatted with `gofmt`
- All tests must pass
- New code should be covered by tests
- Documentation should be updated if necessary
- Benchmark results should not show significant performance regression

## Style Guides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- When only changing documentation, include `[docs]` in the commit title
- Consider starting the commit message with an applicable emoji:
    - 🐛 `:bug:` when fixing a bug
    - ✨ `:sparkles:` when adding a feature
    - 📝 `:memo:` when writing docs
    - 🔥 `:fire:` when removing code or files
    - ⚡️ `:zap:` when improving performance

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format your code with `gofmt`
- Document all exported functions, types, and constants
- Keep function complexity manageable
- Use meaningful variable and function names
- Write clear and comprehensive comments

## Additional Notes

### Issue and Pull Request Labels

This project uses labels to help track and manage issues and pull requests.

- `bug`: Something isn't working
- `documentation`: Improvements or additions to documentation
- `enhancement`: New feature or request
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention is needed
- `performance`: Related to performance issues

## Acknowledgements

This document was adapted from the [Atom Contributing Guide](https://github.com/atom/atom/blob/master/CONTRIBUTING.md) and various other open-source projects.

## Questions?

If you have any questions about contributing, please open an issue with your question.

Thank you for contributing to Crumbs!
