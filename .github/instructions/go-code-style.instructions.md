---
applyTo: "**/*.go"
---

# Go Coding Style

## General Guidelines

- Use as much standard library as possible. Resort to third-party libraries only when necessary.
- Break long lines at 120 characters.
- Always build testable code. Use interfaces to decouple components and make them easier to test.
- Verify structs implement the interfaces they claim to implement by adding `var _ InterfaceName = (*StructName)(nil)`
  just below the struct definition.
- Provide a "new" function for each struct, call it `NewStructName`, and always define it as the first function
  after the struct definition, after the `var _` line (if any).
- Try to space out the code in a way that makes it easy to read.
  For example, if you have a function with multiple arguments, try to align them vertically.
  Insert a blank line between logical sections of code. "Unwrapping" errors to check for nil values is not a logical section,
  and could be considered part of the same section.
