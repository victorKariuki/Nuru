# UI Integration Audit

**Scope:** Agent actions immediately reflected in UI  
**Workspace:** Nuru (Go scripting language with REPL)  
**Date:** 2025-02-22  

---

## 1. How state changes propagate

- **REPL (interactive):** `repl.Start()` uses a single shared `*object.Environment` (`d.env`). Each line is: Lex → Parse → `evaluator.Eval(program, env)` → if result is non-null, `fmt.Println(styles.ReplStyle.Render(evaluated.Inspect()))`. Parser errors are printed before Eval; runtime errors are returned as `*object.Error` and printed via the same non-null branch. No streaming, polling, event bus, or file watching.
- **File run:** `repl.ReadWithEnv(contents, env)` runs the same pipeline once (one Eval for the whole program). Only the **last** statement’s return value is shown; `evalProgram` overwrites `result` each statement. Same stdout-only, synchronous path.
- **Docs TUI:** `repl.Docs()` is a Bubbletea app with a playground that runs code via `evaluator.Eval(program, env)` and sets a viewport’s content to either parser errors or the last non-null result’s `Inspect()`. Same semantics as REPL: one-shot run, last value only, no streaming.
- **HTTP (net package):** Handlers run in Go’s `http.Server` goroutines. “UI” is the HTTP response: handler code can call `jibu.andika(...)`, which writes to `http.ResponseWriter`. That is immediate for the request; there is no REPL or shared REPL state with the server.

**Mechanisms present:** Synchronous return-value path only. No streaming, no polling, no shared event bus, no file watchers. REPL and file run use the same `fmt` + last-result rule.

---

## 2. Agent Action → UI Update Analysis

| Agent Action / State Change              | UI Mechanism                    | Immediate? | Notes |
|-----------------------------------------|----------------------------------|------------|--------|
| **let x = val** (LetStatement)           | None (eval returns `nil`)        | N/A        | Silent: REPL shows nothing; state persists in `env`. |
| **x = val** (Assign)                     | Last-result print               | Yes        | `env.Set` returns value; REPL prints it. |
| **x += 1** etc (AssignEqual)             | Last-result print               | Yes        | Returns new value; REPL prints it. |
| **ingiza X** / import (Mapper or file)   | None (returns `NULL`)            | N/A        | Silent: bindings added to `env`, no REPL output. |
| **fanya pakeji** / package def           | Last-result print               | Yes        | Returns package object; REPL prints. |
| **unda fn () { }** / function literal    | Last-result print               | Yes        | Expression returns function; REPL prints. |
| **andika(...)**                          | Builtin → `fmt.Println`          | Yes        | Direct stdout; not the “last value” path. |
| **jaza(...)**                            | Return value if used as expr    | Yes        | Builtin returns string; shown if line is expression. |
| **obj.x = val** (PropertyAssignment)     | None (returns `NULL`)           | N/A        | Silent: instance/dict/package updated, no output. |
| **arr[i] = x** / **dict[k] = v**         | Last-result print               | Yes        | Evaluator returns assigned value. |
| **i++** / postfix                        | Last-result print               | Yes        | `env.Set` returns new value; REPL prints. |
| Parser errors                            | `fmt.Println` before Eval       | Yes        | Shown; Eval still runs (possible partial AST). |
| Runtime errors (*object.Error)           | Last-result print               | Yes        | `Type() != NULL_OBJ`, so `Inspect()` printed. |
| Block / multi-statement program          | Last statement only             | Yes (last) | Intermediate statements (let, import, etc.) not shown. |
| HTTP handler (e.g. **jibu.andika(...)**)| `w.Write` / response            | Yes        | Per-request; no REPL. |
| VerboseGC / PrintObjectStats            | Opt-in `fmt` in object pkg      | Yes        | Only when `--verbose` or explicit debug; not main UI. |

---

## 3. Score: ~50–55% (percentage of state-changing actions that surface in UI)

- **Reflected:** Assign, AssignEqual, package def, function def, andika, jaza (as expression), index assignment, postfix, parser/runtime errors, last expression in block/file, HTTP response write.
- **Silent:** Let (variable binding), import (Mapper or file), property assignment (obj.x = val), and any statement that returns `nil`/`NULL` and is not the last in a block (e.g. intermediate lets in a file).

Rough count: a large share of **bindings** (let, import, property assign) do not produce any UI output; **expressions** and **assignments that return a value** do.

---

## 4. Silent actions (anti-patterns)

- **LetStatement:** `env.Set` then no return → evaluator returns `nil` → REPL shows nothing. User has no confirmation that the binding was created except by typing the variable on the next line.
- **Import (ingiza):** Both `evalImport` (Mapper) and `importFile` end with `return NULL` after `env.Set`. No “imported X” or module summary.
- **PropertyAssignment:** `evalPropertyAssignment` always returns `NULL` after `obj.Env.Set`/`dict.Set`. No feedback for `obj.x = val` or `pakeji.something = val`.
- **File / block:** Only the last statement’s value is shown. All previous statements (including lets, imports, property assigns) are silent; multi-statement files look like “no output” until the last expression.
- **Parser errors then Eval:** REPL and ReadWithEnv still call `evaluator.Eval(program, env)` after printing parser errors. Partial/broken AST may be evaluated; no extra guard.

---

## 5. Recommendations

1. **Optionally echo bindings in REPL:** For `let x = val` and successful `ingiza X`, consider a small line (e.g. “x = …” or “Imeload: X”) so state changes are visible. Could be behind a flag (e.g. `--echo-bindings`).
2. **Return assigned value from property assignment:** Have `evalPropertyAssignment` return `val` (or the updated object) instead of `NULL` so `obj.x = 5` shows a result like other assignments.
3. **Import return value:** Consider returning the module (or a small “imported” object) from `evalImport`/`importFile` so REPL can show “module X” and imports are not silent.
4. **File run:** Consider “verbose” mode that prints each statement’s result (or at least binding echoes) so file execution is not only “last value.”
5. **Docs playground:** Align with REPL: same rules for what is shown; consider echoing binds/imports if you add that to the REPL.
6. **Parser errors:** Consider skipping `Eval` when `len(p.Errors()) != 0` to avoid partial evaluation and clarify that “only fixed code runs.”

---

*Audit performed by tracing: main.go → repl.Start/ReadWithEnv → dummy.executor / Read/ReadWithEnv → evaluator.Eval; evaluator cases for Let, Assign, AssignEqual, Import, Package, PropertyAssignment, Block, and builtins (andika, jaza); object.Environment.Set return value; module HTTP handler and EvaluatorCallback.*
