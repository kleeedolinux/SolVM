# Contributing to SolVM

First off, thank you for considering contributing to SolVM! We're excited to see the community get involved and help make SolVM a more powerful and robust runtime for Lua. Your contributions, whether big or small, are valuable.

This document provides guidelines to help you contribute effectively.

## How Can I Contribute?

There are many ways to contribute to SolVM:

*   **Reporting Bugs:** If you find a bug, please let us know!
*   **Suggesting Enhancements:** Have an idea for a new feature or an improvement to an existing one? We'd love to hear it.
*   **Writing Code:** If you're up for it, you can contribute by fixing bugs or implementing new features.
*   **Improving Documentation:** Clear documentation is crucial. If you see areas for improvement or find something confusing, please help us make it better.
*   **Writing Examples:** More examples of how to use SolVM's features are always welcome.
*   **Spreading the Word:** Tell your friends and colleagues about SolVM!

## Getting Started

Before you start, please:

1.  **Familiarize yourself with the project:** Read the `README.md` and the `DOC.md` to understand SolVM's goals, architecture, and existing features.
2.  **Set up your development environment:**
    *   You'll need Go (check `go.mod` for the recommended version).
    *   Clone the repository: `git clone https://github.com/kleeedolinux/SolVM.git`
    *   Navigate into the project directory: `cd SolVM`
    *   Ensure you can build the project (e.g., `go build -o solvm main.go` or use the provided build scripts).
3.  **Check existing issues and discussions:** Someone might have already reported the bug or suggested the feature you have in mind.

## Reporting Bugs

A good bug report helps us identify and fix issues faster. When reporting a bug, please include:

*   **A clear and descriptive title.**
*   **The version of SolVM you are using.** (You can usually get this by running `solvm --version` if implemented, or state the commit hash).
*   **Your operating system and version.**
*   **Steps to reproduce the bug:** Be as specific as possible. Provide a minimal Lua script that demonstrates the issue.
*   **What you expected to happen.**
*   **What actually happened:** Include any error messages or logs.
*   **If possible, a minimal, reproducible example.** This is incredibly helpful.

**Security Vulnerabilities:**
**Do not open a public GitHub issue for security vulnerabilities.** Instead, please email me@juliaklee.wtf with the details. We take security seriously and will address any vulnerabilities promptly.

## Suggesting Enhancements or New Features

We love new ideas! When suggesting an enhancement:

*   **Explain the problem you're trying to solve** or the use case the new feature would address.
*   **Describe your proposed solution** in as much detail as possible.
*   **Explain why this enhancement would be useful** to SolVM users.
*   Consider discussing your idea in an [Issue](https://github.com/kleeedolinux/SolVM/issues) first, especially for larger changes, to gather feedback before investing too much time in implementation.

## Submitting Pull Requests (PRs)

Ready to contribute code? Great! Here's how to make the process smoother:

1.  **Fork the repository** on GitHub.
2.  **Create a new branch** for your changes: `git checkout -b feature/your-awesome-feature` or `fix/bug-description`.
3.  **Make your changes:**
    *   **Follow existing code style:** Try to maintain consistency with the existing codebase (Go for the runtime, Lua for examples/tests where applicable). We generally follow standard Go formatting (`gofmt`).
    *   **Write clear and concise commit messages:** See "Commit Message Guidelines" below.
    *   **Add tests for your changes:** If you're fixing a bug, include a test that reproduces the bug and is fixed by your patch. If you're adding a new feature, include tests that cover its functionality.
    *   **Update documentation:** If your change affects user-facing behavior or adds new features, please update `DOC.md` or other relevant documentation.
4.  **Ensure your changes build successfully** and all tests pass.
5.  **Push your branch** to your fork: `git push origin feature/your-awesome-feature`.
6.  **Open a Pull Request** against the `main` branch of the `kleeedolinux/SolVM` repository.
    *   Provide a **clear title and description** for your PR. Explain the "what" and "why" of your changes.
    *   **Link to any relevant issues** (e.g., "Fixes #123").
    *   Be prepared to discuss your changes and make adjustments based on feedback from maintainers.

**A Note on Cosmetic Changes:**
Pull requests that only involve whitespace changes, code reformatting, or other purely cosmetic modifications (without affecting functionality, stability, or testability) might not be accepted. We want to keep the commit history focused on meaningful changes.

## Commit Message Guidelines

Clear commit messages help everyone understand the history of the project. We generally follow these conventions:

*   **Use the present tense and imperative mood:** "Fix bug in HTTP server" instead of "Fixed bug..." or "Fixes bug...".
*   **Keep the subject line short (around 50 characters).**
*   **If more explanation is needed, provide a more detailed body after a blank line.**
*   **Reference relevant issue numbers** in the body if applicable (e.g., `Fixes #123`).

Example:
```
feat: Add support for YAML parsing in 'data' module

This commit introduces `yaml.encode` and `yaml.decode` functions
to the built-in 'data' module, allowing users to easily work
with YAML-formatted data within their SolVM scripts.

Includes unit tests for both encoding and decoding various
YAML structures. Updates DOC.md with the new functions.

Closes #456
```

## Code of Conduct

SolVM is dedicated to providing a welcoming and inclusive environment for everyone. All contributors are expected to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md) (you'll need to create this file, often based on a template like the Contributor Covenant). Please be respectful and constructive in all interactions.

## Questions?

If you have questions about how to contribute, how SolVM works, or anything else, feel free to:
*   Open an [Issue](https://github.com/kleeedolinux/SolVM/issues) for discussion.

---

Thank you for taking the time to contribute to SolVM. We appreciate your help in building a better runtime for the Lua community!
