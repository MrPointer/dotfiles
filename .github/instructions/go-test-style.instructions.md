---
applyTo: "**/*.go"
---

# Go Test Style

## General Guidelines

- Place tests in the test package of the package being tested.
  - For example, if the package is `lib`, the test package should be named `lib_test`.
- Use the `testify` library for all testing in Go.
  - Always use the `require` package from `testify` when checking the `error` type.
- Each test should verify a single behavior or property. Do not test multiple behaviors in a single test.
- Name tests based on their behavior, using descriptive and natural language.
  - Example: `Test_CompatibilityConfigCanBeLoadedFromFile` checks if the `CompatibilityConfig` struct can be loaded from a file.
  - Test names should describe what the test does and what it verifies, not implementation details.
  - Use the `Test_` prefix for test functions.
  - If a condition is crucial to the test, include it in the test name.
    - Example: `Test_CompatibilityConfigCanBeLoadedFromFile_WhenFileExists` indicates that the test checks loading from an existing file.
    - Separate different conditions with underscores.
- Use table-driven tests when appropriate. Table-driven tests define a set of inputs and expected outputs in a table, and iterate over the table to run the tests. This pattern makes it easy to add new test cases and keeps the code clean and maintainable.

## Types of Tests

### Unit Tests

- Unit tests verify a single function or method in isolation.
- Use mocks to isolate the function or method being tested.
  - Use the [moq](https://github.com/matryer/moq) package to generate mocks.

### Integration Tests

- Integration tests verify the interaction between multiple functions or methods.
- Integration tests also cover OS-dependent interactions (anything beyond CPU and memory).
- For every integration test, allow opting out by using `testing.Short()`.
  - This is useful for running only unit tests in CI/CD pipelines.

## Testing Tech Stack

- [testify]: A Go library for writing tests and assertions.
- [moq]: A Go library for generating mocks for testing.
- [mockery]: A Go library for generating mock objects. It is used to generate mock objects in the project.

## Using Mocks

- Use [mockery] to generate mocks for interfaces. Run the command `mockery` (with no arguments) in the root directory of the Go module (for example, the `go-port` directory).
- In test code, use the generated mocks to test your code. The generated mocks are compatible with the `moq` library, so you can use `moq` features in your tests.

[testify]: https://github.com/stretchr/testify
[moq]: https://github.com/matryer/moq
[mockery]: https://github.com/vektra/mockery
