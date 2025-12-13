---
applyTo: "**/*.go"
---

- Always return wrapped errors with context using fmt.Errorf("context: %w", err).
- Use testify require for assertions in tests.
- Use defer func() { _ = someFunc() }() style to prevent linter warnings for ignored errors.