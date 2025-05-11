---
applyTo: "**/*_test.go"
---

# Go Test Style

## General Guidelines

- Don't use `testify` for testing. Use the standard library `testing` package instead.
- Always test a single thing in a test. Don't test multiple things in a single test.
- Name tests based on their behavior, expressed in natural language.
  For example, `TestCompatibilityConfigCanBeLoadedFromFile`, which checks if the `CompatibilityConfig`
  struct can be loaded from a file, is a good name for a test.
  It describes what the test does and what it tests, but doesn't necessarily focus on the implementation
  or technical details.
- Use table-driven tests where appropriate. Table-driven tests are a common pattern in Go testing.
  They allow you to define a set of inputs and expected outputs in a table, and then iterate over
  the table to run the tests. This makes it easy to add new test cases and keeps the code clean.

## Different Types of Tests

### Unit Tests

- Unit tests are tests that test a single function or method in isolation.
- Use mocks to isolate the function or method being tested.
  Use the [moq][moq] package to generate mocks.

### Integration Tests

- Integration tests are tests that test multiple functions or methods together.
  Integration tests also test OS-dependent interactions, which is anything but CPU and memory.
- For every integration test, make it possible to opt-out of the test by using `testing.Short()`.
  This is useful for running tests in CI/CD pipelines where you want to run only unit tests.

### System Tests

- System tests are tests that test the entire system, or that test the system in a specific environment.
- System tests are usually run in a separate environment, such as a Docker container or a virtual machine.
- Such tests should use [testcontainers-go][testcontainers-go] to run the tests in a container.
  Use it for any test that requires a specific environment or setup, or one that could ruin
  the host system if it were to run on it. For example, a test to check that homebrew can be installed
  should run in a container, as it could modify the host system if it were to run on it.

[moq]: https://github.com/matryer/moq
[testcontainers-go]: https://golang.testcontainers.org/
