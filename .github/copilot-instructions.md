---
applyTo: "**/*.go"
---

- Always return wrapped errors with context using fmt.Errorf("context: %w", err).
- Use testify require for assertions in tests.