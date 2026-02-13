# Interactive UI Reference

Interactive UI is built with [Huh](https://github.com/charmbracelet/huh) for terminal-based forms and selection prompts. All interactive components live in `cli/`.

## Architecture Pattern

Every interactive UI component follows this layering:

1. **Define an interface** for the interaction (enables testing with mocks)
2. **Implement with Huh** as the concrete implementation
3. **Provide a `NewDefault*` constructor** that wires up Huh internally
4. **Consume via dependency injection** in command logic

See `cli/selector.go` for the generic `Selector[T]` interface and its Huh implementation. See `cli/gpg_selector.go` for a domain-specific wrapper example.

## Key Behaviors

- **Single-item optimization**: if only one option exists, return it directly without prompting
- **Empty-list error**: return an error if no items are available for selection
- **Label mapping**: use `SelectWithLabels` when display text differs from the underlying value

## Non-Interactive Fallback

When `--non-interactive` is set, skip all Huh prompts entirely. Use flag values or defaults instead, and log what was skipped as a warning. Check `IsNonInteractive()` before any interactive operation.

## Progress Transitions

Use `StartInteractiveProgress` / `FinishInteractiveProgress` / `FailInteractiveProgress` when transitioning between progress display and interactive UI. This properly pauses spinners during user input.

## Key Rules

- Never call Huh directly from commands — always go through an interface in `cli/`
- Use generics for reusable selection patterns (see existing `Selector[T]`)
- Provide `NewDefault*` constructors that wire Huh internally
- Respect `--non-interactive` in all UI components
