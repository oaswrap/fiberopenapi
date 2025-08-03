# 📦 [DEPRECATED] `oaswrap/fiberopenapi`

> **⚠️ This repository is archived.**

The `oaswrap/echoopenapi` adapter has been **migrated to the official mono-repo**:

👉 **New location:** [`github.com/oaswrap/spec/adapters/fiberopenapi`](https://github.com/oaswrap/spec)

## 📚 Why was this moved?

To simplify development and versioning, **`oaswrap`** now uses a **monorepo**:
- ✅ The **core OpenAPI spec generator** and **framework adapters** (Gin, Fiber, Echo, etc.) live together.
- ✅ All adapters share the same version and commit.
- ✅ Easier contributions and issue tracking.

## 🚀 How to use the new adapter

In your `go.mod`:
```go
require (
    github.com/oaswrap/spec v0.2.0
    github.com/oaswrap/spec/adapters/fiberopenapi v0.0.0-<pseudo-version>
)
```
