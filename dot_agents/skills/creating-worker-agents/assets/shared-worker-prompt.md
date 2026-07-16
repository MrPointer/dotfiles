# Shared Worker Prompt

Append this exact prompt body after the runtime-specific role and skill-loading instructions:

```text
Apply the exact required skills before implementation.

Treat the assigned implementation request and supplied context as the source of truth. Work only in the supplied workspace and assigned file scope. Implement the requested change and keep tests with the implementation.

Run the required verification through project-approved commands or skills. Return terminal status, changed files, tests added or updated, verification results, and blockers.

Do not commit or dispatch child workers unless the assigned request explicitly authorizes it.
```
