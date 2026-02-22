# Nuru LSP and DAP

This document describes the **Language Server Protocol (LSP)** and **Debug Adapter Protocol (DAP)** implementations for Nuru: what they provide, how to build and run them, and how they integrate with the VS Code extension.

---

## Overview

| Component   | Binary      | Role |
|------------|-------------|------|
| **LSP**    | `nuru-lsp`  | Editor support: diagnostics, go-to-definition, hover, completion, outline, rename, code actions, formatting. |
| **DAP**    | `nuru-dap`  | Debugging: launch script, breakpoints, call stack, scopes, variables, continue/step. |

Both servers speak JSON-RPC over stdio. The VS Code extension starts them when you work with `.nr` / `.sw` files (LSP) or start a debug session (DAP).

---

## Language Server (LSP)

### Build and run

From the repository root:

```bash
go build -o nuru-lsp ./cmd/nuru-lsp
```

Run manually (stdio):

```bash
./nuru-lsp
```

The extension normally starts the server automatically; you only need the binary on your PATH or set **Nuru: Language Server Path** in settings.

### Capabilities

The server advertises these LSP capabilities:

| Capability | Description |
|------------|-------------|
| **Text document sync** | Full document sync on open/change/close. |
| **Diagnostics** | Parser errors (and analysis errors) as diagnostics; published on open, change, and close. |
| **Go to definition** | Jump to the definition of an identifier (variable, function) using the symbol table. |
| **Hover** | Hover over an identifier to see type/kind and definition location. |
| **Completion** | Trigger completion (e.g. Ctrl+Space) for identifiers in scope. |
| **Document symbol** | Outline view: list of top-level symbols (variables, functions) in the file. |
| **Rename** | Rename symbol at cursor; all references in the document are updated. |
| **Code actions** | Quick fixes (e.g. insert missing `;`) and “Organize imports” (sort and dedupe `tumia` lines). |
| **Formatting** | Format entire document (e.g. Format Document). Uses the `format` package (AST-based). |

### Implementation notes

- **Parser errors** are converted to LSP diagnostics; messages like `Mstari N: ...` are parsed to get the line (1-based in Nuru, 0-based in LSP).
- **Symbol table** (`lsp/symbols`) is built from the AST per document; it drives definition, hover, completion, document symbol, and rename.
- **Rename** uses reference finding over the AST (`lsp/references.go`) to compute all edits in the same file.
- **Formatting** is implemented in the `format` package: it rewrites the AST to a consistent style (indent, semicolons, line breaks). The LSP server calls `format.Program` with options from the client (tab size, insert spaces).

### Source layout

- `cmd/nuru-lsp/main.go` — entrypoint; runs `lsp.NewServer().Run(os.Stdin, os.Stdout)`.
- `lsp/server.go` — request routing and handlers (initialize, didOpen/didChange/didClose, definition, hover, completion, documentSymbol, rename, codeAction, formatting).
- `lsp/protocol.go` — LSP types (Position, Range, Diagnostic, params, capabilities).
- `lsp/diagnostics.go` — parser/analysis errors to diagnostics; diagnostic codes for code actions.
- `lsp/position.go` — conversion between LSP positions and AST/lexer positions.
- `lsp/references.go` — find all references to an identifier in the program (for rename).
- `lsp/symbols/table.go` — symbol table built from AST (definitions, scopes).
- `format/format.go` — AST-based formatter used by the LSP formatting handler.

---

## Debug Adapter (DAP)

### Build and run

From the repository root:

```bash
go build -o nuru-dap ./cmd/nuru-dap
```

The debug adapter is intended to be launched by the IDE (e.g. VS Code). For a manual test over stdio you can run `./nuru-dap` and send DAP JSON-RPC messages; normally the extension starts it with a launch config.

### Capabilities

| DAP request/event | Support |
|-------------------|--------|
| **initialize** | Supported; sends `initialized` event. |
| **launch** | Launches the Nuru script in `program`; parses and runs with the evaluator. |
| **setBreakpoints** | Sets breakpoints by source path and line (DAP: 0-based line; internally converted to 1-based). |
| **configurationDone** | Acknowledged. |
| **continue** | Resumes execution (unblocks the debug hook). |
| **next / stepIn / stepOut** | Treated as continue (single-threaded; step semantics not yet differentiated). |
| **disconnect** | Stops the session and exits the adapter. |
| **threads** | Returns a single thread (id: 1, name: "main"). |
| **stackTrace** | Returns the current call stack (frames with id, name, line, source path). |
| **scopes** | Returns “Locals” scope for the selected frame (variablesReference = frameId + 1000). |
| **variables** | Returns variables for a scope (by variablesReference); for locals, reads from the frame’s environment. |
| **stopped** (event) | Fired when a breakpoint is hit (reason: breakpoint, threadId: 1, line: 0-based). |
| **terminated** (event) | Fired when the program finishes. |

### Implementation notes

- **Execution** runs in a goroutine; the main loop handles DAP requests. When a breakpoint is hit, the evaluator’s `DebugHook` runs and blocks until the client sends continue/step/disconnect.
- **Breakpoints** are stored per source path and 1-based line; the hook checks the current node’s line against this set.
- **Call stack** is maintained via `evaluator.DebugStackPush` / `DebugStackPop`; each frame stores name, line, path, and environment. The evaluator calls these when entering/leaving blocks and function calls.
- **Variables** are read from `object.Environment` for the selected frame; names come from `env.Names()`, values from `Inspect()`.
- **Line numbers**: DAP uses 0-based lines in requests/responses; the adapter converts to 1-based for the evaluator and back to 0-based for stack/stopped events.

### Source layout

- `cmd/nuru-dap/main.go` — DAP server: readMessage/send loop, request dispatch, launch, setBreakpoints, continue/step, threads, stackTrace, scopes, variables; integration with evaluator hooks.
- `evaluator/evaluator.go` — `DebugHook`, `DebugStackPush`, `DebugStackPop`; hook is invoked at each evaluator step; stack push/pop used for frames.
- `evaluator/block.go` — stack push/pop for block execution.
- `object/environment.go` — `Names()` for listing variable names in a scope (used by DAP variables).

---

## VS Code extension

- **LSP**: The extension starts `nuru-lsp` when the language is activated (e.g. opening a `.nr` file). Settings: `nuru.enableLanguageServer`, `nuru.languageServerPath`.
- **DAP**: The extension registers a debug adapter descriptor factory for type `nuru`; when you start a debug session with a “Launch Nuru file” config, it runs `nuru-dap` and passes `program` (e.g. `${file}`). Setting: `nuru.debugAdapterPath`.

See [extensions/vscode/README.md](../../extensions/vscode/README.md) for user-facing setup (build, PATH, settings).

---

## Tests and CI

- **LSP**: `go test ./lsp/... ./lsp/symbols/...`
- **Format**: `go test ./format/...`
- **Build**: `go build ./cmd/nuru-lsp/` and `go build ./cmd/nuru-dap/`

These are included in `make test`.
