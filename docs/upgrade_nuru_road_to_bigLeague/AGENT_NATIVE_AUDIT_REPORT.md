# Agent-Native Architecture Review: Nuru

**Project:** Nuru — A Swahili programming language (interpreter + REPL in Go).  
**Audit date:** 2026-02-22.  
**Scope:** Full codebase audit against 8 agent-native principles.

**Note:** Nuru is a scripting language and runtime, not an application with an embedded AI agent. There is no MCP server, no agent tools, and no system prompt. Scores reflect the current state; many “gaps” are expected until/unless agent integration is added.

---

## Overall Score Summary

| Core Principle           | Score    | Percentage | Status |
|--------------------------|----------|------------|--------|
| Action Parity            | 0/10     | 0%         | ❌     |
| Tools as Primitives      | 161/164  | 98%        | ✅     |
| Context Injection        | 0/6      | 0%         | ❌     |
| Shared Workspace         | 6/6      | 100%       | ✅     |
| CRUD Completeness        | 2/4      | 50%        | ⚠️     |
| UI Integration           | ~12/22   | ~55%       | ⚠️     |
| Capability Discovery     | 2/7      | 29%        | ❌     |
| Prompt-Native Features   | 0/N      | 0%         | ❌     |

**Overall Agent-Native Score: ~41%**

### Status Legend

- ✅ Excellent (80%+)
- ⚠️ Partial (50–79%)
- ❌ Needs Work (&lt;50%)

---

## 1. Action Parity — 0/10 (0%)

**Finding:** No agent exists. All user actions (REPL, run script, docs TUI, analyze, visualize deps, help, version) are only available via the CLI/REPL; there are no MCP or agent tools.

| Action                    | Location              | Agent Tool | Status   |
|---------------------------|-----------------------|------------|----------|
| Start REPL                | main.go               | None       | Missing  |
| Run script file           | main.go               | None       | Missing  |
| CLI help / version / docs | main.go               | None       | Missing  |
| Analyze (--chambua)       | analysis/cli.go       | None       | Missing  |
| Visualize deps (--tengeneza) | analysis/cli.go    | None       | Missing  |
| REPL evaluate / exit      | repl/repl.go          | None       | Missing  |

**Recommendation:** Add an MCP server (or equivalent) exposing tools: run script, evaluate snippet, analyze, visualize deps, return doc content, version/help — or a single “run Nuru CLI” tool with args for parity.

---

## 2. Tools as Primitives — 161/164 (98%)

**Finding:** The “tool” surface is the language itself (global builtins + stdlib modules). Almost all are single-capability primitives. Three combined helpers:

| Tool              | Type     | Reasoning                          |
|-------------------|----------|------------------------------------|
| jsoni.soma        | WORKFLOW | Read file + decode (I/O + parse)   |
| jsoni.hifadhi     | WORKFLOW | Encode + write file                |
| http.sajiliNjia   | WORKFLOW | Register route + handler           |

**Recommendation:** Keep design; document primitive composition (e.g. `faili.soma` + `jsoni.dikodi`). If an agent layer is added, expose primitives there.

---

## 3. Context Injection — 0/6 (0%)

**Finding:** No agent or system prompt. No code injects “context” into an AI. REPL uses a fixed prompt prefix; `__FILE__`/`__DIR__` are set for file runs but not fed into any prompt.

**Recommendation:** When/if adding an agent, introduce a single system-prompt builder and inject: available resources (builtins, modules), workspace state (`__FILE__`, `__DIR__`), session bindings (REPL env), capabilities list.

---

## 4. Shared Workspace — 6/6 (100%)

**Finding:** Single actor (script author). One `Environment` per run/REPL, same file I/O, same stdio, no separate “agent” data space. No sandbox isolation anti-pattern.

**Recommendation:** If an agent is added, keep one workspace: agent and user share the same env and paths.

---

## 5. CRUD Completeness — 2/4 (50%)

**Finding:** Among mutable, collection-like entities (Array, Dict, Set, File), only **Set** and **File** have full CRUD. **Array** and **Dict** lack Delete in the language (Go has `Remove`/`Delete` but no exposed builtin/method).

**Recommendation:** Expose delete for Array (e.g. `ondoaKiashiria(i)` or `ondoaThamani(x)`) and Dict (e.g. `ondoa(ufunguo)` or `futa(ufunguo)`). Document immutability of value types so “no Update/Delete” is explicit.

---

## 6. UI Integration — ~55%

**Finding:** REPL shows only the last expression value. Silent state changes: `let x = val`, `ingiza X`, `obj.x = val`, and non-last statements in a block. No streaming, polling, or event bus.

**Recommendation:** Optional binding echo for `let` and `ingiza`; have property assignment return the value; consider optional per-statement or binding echo for file runs; skip Eval when there are parser errors.

---

## 7. Capability Discovery — 2/7 (29%)

**Finding:** Help docs (TUI + embedded docs) and CLI help exist. Missing: onboarding, REPL capability hints, self-description (e.g. `help(andika)`), suggested prompts, richer empty-state guidance, slash commands (e.g. `/help`, `/docs`).

**Recommendation:** Add REPL empty-state line (“Try: andika('habari')”, “Docs: nuru --nyaraka”); in-REPL `/help` or `msaada`; completer with builtins/keywords; optional `Help` on builtins and `msaada(name)`.

---

## 8. Prompt-Native Features — 0%

**Finding:** All behavior is code-defined in Go (builtins, modules, methods, semantics, errors, REPL strings). POLICY.md specifies design outcomes in prose but is not read at runtime.

**Recommendation:** For a classic interpreter this is appropriate. If you want some prompt-native behavior later: move user-facing strings to data files; optionally make POLICY machine-readable for tooling/CI.

---

## Top 10 Recommendations by Impact

| Priority | Action | Principle           | Effort  |
|----------|--------|---------------------|---------|
| 1        | Add MCP server with run script, evaluate, analyze, deps, docs (or one “run CLI” tool) | Action Parity | High    |
| 2        | Inject dynamic context when adding an agent (resources, workspace, session, capabilities) | Context Injection | Medium  |
| 3        | Expose Array/Dict delete in the language (ondoa/futa) | CRUD Completeness | Low     |
| 4        | REPL: echo bindings for `let`/`ingiza`; property assign returns value | UI Integration | Low     |
| 5        | REPL empty-state hint + in-REPL `/help` or `msaada` | Capability Discovery | Low     |
| 6        | Completer with builtins/keywords; optional `msaada(name)` for builtins | Capability Discovery | Medium  |
| 7        | Optional “run Nuru CLI” agent tool for quick parity without full MCP | Action Parity | Medium  |
| 8        | Document primitive composition for jsoni (faili.soma + jsoni.dikodi, etc.) | Tools as Primitives | Low     |
| 9        | Skip Eval when parser has errors to avoid partial run | UI Integration | Low     |
| 10       | Document single-workspace model and “no agent yet” in architecture docs | Shared Workspace / General | Low     |

---

## What’s Working Well

1. **Tools as Primitives (98%)** — Stdlib and builtins are largely single-capability; only three combined helpers, and the language is a solid base for future agent tools.
2. **Shared Workspace (100%)** — Single data space per run, no agent/user split, no sandbox isolation.
3. **Rich, bilingual docs** — `repl/docs/` and docs TUI (`--nyaraka`) give strong capability discovery outside the REPL.
4. **Clear CLI surface** — One entry point, clear flags (run, analyze, visualize, docs, help, version), easy to map to future agent tools.
5. **CRUD in Go for collections** — Array and Dict already have Remove/Delete in Go; exposing them in the language is a small, high-impact change.

---

## Success Criteria Checklist

- [ ] Agent can achieve what users can via CLI/REPL (parity) — **blocked: no agent**
- [x] Tools are atomic primitives; combined helpers documented — **language layer: yes**
- [ ] New features addable via prompts — **N/A: no agent**
- [ ] Agent handles unanticipated requests — **N/A: no agent**
- [x] Single workspace; no isolated agent data — **yes**
- [ ] System prompt includes dynamic context — **no agent**
- [ ] Every entity has full CRUD where applicable — **Array/Dict missing Delete**
- [ ] State changes visible in UI (REPL) — **partial: let/import/property silent**
- [ ] Users can discover capabilities — **partial: docs good; REPL discovery weak**

---

*Report generated by agent-native-audit (8 parallel sub-agents). Detailed per-principle write-ups may exist in `docs/` (e.g. TOOLS_AS_PRIMITIVES_AUDIT.md, CRUD_COMPLETENESS_AUDIT.md, UI_INTEGRATION_AUDIT.md).*
