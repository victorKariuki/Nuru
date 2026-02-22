# Curriculum: Units the codebase can support

A **computer-science style progression**: environment → values and expressions → control flow → abstraction → data structures → I/O and data → text and time → numerics → networks → system. Each unit is teachable today using existing reference docs and/or runnable examples.

**Sources:** reference docs in `repl/docs/en/`, scripts in `examples/`, and language design in `docs/POLICY.md`.

---

## Part 0 — Environment and first program

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 0.1 | **Getting started** | Tooling, REPL, run script | builtins, comments, identifiers, keywords, README | example.nr, smoke_no_io.nr | — |

*Learning focus:* Install Nuru, run REPL, run a file, use `andika`/`jaza`, comments and naming (identifiers, keywords).

---

## Part 1 — Values, types, and expressions

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 1.1 | **Variables and types** | Binding, type of value | identifiers, numbers, strings, bool, null, builtins | example.nr | 0.1 |
| 1.2 | **Operators** | Expressions, precedence | operators | example.nr | 1.1 |
| 1.3 | **Numbers (primitives)** | Integers, floats, literals | numbers | example.nr | 1.1 |
| 1.4 | **Strings (primitives)** | Literals, concatenation | strings | example.nr | 1.1 |
| 1.5 | **Booleans and null** | Truth, absence of value | bool, null | example.nr | 1.1 |

*Learning focus:* `fanya`, `aina`, `badilisha`, `namba`, `tungo`; arithmetic and comparison; string basics; boolean logic and `tupu`.

---

## Part 2 — Control flow

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 2.1 | **Conditionals** | Branching | ifStatements | conditional_stress.nr | 1.x |
| 2.2 | **Loops: for and while** | Iteration, termination | for, while | for_in_stress.nr, while_loop_stress.nr, control_flow_combined_stress.nr, c_style_for_example.nr, simple_for_in_test.nr | 2.1 |
| 2.3 | **Switch** | Multi-way dispatch | switch | — | 2.1 |

*Learning focus:* `kama` / `au kama` / `sivyo`; `kwa` / `wakati`; `vunja` / `endelea`; C-style for; `badili` / `kawaida`.

---

## Part 3 — Abstraction

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 3.1 | **Functions** | Procedures, parameters, return | function | example.nr | 2.x |
| 3.2 | **Recursion and closures** | Self-reference, capture | function | example.nr | 3.1 |
| 3.3 | **Packages** | Modules, reuse | packages | — | 3.1 |

*Learning focus:* `unda`, parameters, default args, `rudisha`; recursion; closures; `tumia <pakeji>`.

---

## Part 4 — Data structures

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 4.1 | **Arrays (sequences)** | Ordered collections, indexing | arrays | example.nr, reduce.nr, sorting_algorithm.nr | 3.1 |
| 4.2 | **Strings (sequence methods)** | Length, split, format | strings | example.nr | 4.1 |
| 4.3 | **Range, sets, tuples** | mfululizo, seta, jozi | range, sets, builtins | — | 4.1 |
| 4.4 | **Dictionaries** | Key–value, mapping | dictionaries | example.nr, iterator_example.nr | 4.1 |
| 4.5 | **Iteration (kitanzi)** | Traversal abstraction | for, arrays, dictionaries | iterator_example.nr | 4.1, 4.4 |

*Learning focus:* ORODHA (sukuma, chuja, ramani, unga, …); NENO methods (idadi, gawa, panga, …); mfululizo/seta/jozi; KAMUSI (fungua, funguo, maana); `kitanzi()` and `kwa` over collections.

---

## Part 5 — I/O and data formats

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 5.1 | **Files and filesystem** | Streams, paths on disk | files | fs_full_example.nr, fs_path_json_example.nr, script_relative_path_example.nr | 4.x |
| 5.2 | **Script context** | __FILE__, __DIR__ | — | script_relative_path_example.nr | 5.1 |
| 5.3 | **JSON** | Serialization, structured data | json | jsoni_full_test.nr, fs_path_json_example.nr | 4.4, 5.1 |

*Learning focus:* `tumia faili` (fungua, soma, andika, orodha, tengenezaSarafa, …); script-relative paths; `tumia jsoni` (dikodi, enkodi, soma, hifadhi).

---

## Part 6 — Text and time

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 6.1 | **Regular expressions** | Pattern matching, search/replace | regex | regex_example.nr | 4.2 |
| 6.2 | **Time** | Timestamps, duration, sleep | time | — | 3.1 |

*Learning focus:* `tumia re` (linganisha, tafuta, tafutaZote, vikundi, badilisha, gawa, tayari); `tumia muda` (hasahivi, lala, tangu, ongeza, …).

---

## Part 7 — Numerics

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 7.1 | **Math library (hisabati)** | Constants, sqrt, random, trig | hisabati | big_integer_example.nr | 1.3 |
| 7.2 | **Big integers** | Arbitrary-precision arithmetic | builtins (badilisha), hisabati | big_integer_example.nr | 7.1 |

*Learning focus:* hisabati constants and functions; `badilisha(x, "NAMBA_KUBWA")`, hisabati.namba_kubwa.

---

## Part 8 — Networks and the web

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 8.1 | **HTTP client** | Request/response | net | http_full_example.nr | 4.x |
| 8.2 | **HTTP server** | Handlers, routing | POLICY (http) | http_server_example.nr | 3.1, 8.1 |
| 8.3 | **URLs** | Parsing, building | POLICY (url) | url_full_example.nr | 8.1 |

*Learning focus:* `tumia mtandao` (peruzi, tuma); `tumia http` (sajiliNjia, undaServer, …); `tumia url` (changanua, tengeneza, tatua, …).

---

## Part 9 — System and environment

| Unit | Topic | Concepts | Reference | Examples | Prereqs |
|------|--------|----------|-----------|----------|---------|
| 9.1 | **Paths (njia)** | Path manipulation, glob | POLICY (njia) | njia_glob_example.nr, njia_full_example.nr | 5.1 |
| 9.2 | **Crypto and encoding** | Hashing, Base64 | POLICY (crypto) | crypto_example.nr, simple_base64.nr, advanced_base64_examples.nr | 4.1 |
| 9.3 | **OS (process)** | Exit, exec | POLICY (os) | — | 3.1 |
| 9.4 | **Memory / runtime (mfumo)** | GC, stats, weak refs | POLICY (mfumo) | memory_test.nr | 4.x |

*Learning focus:* `tumia njia` (unganisha, jina, sarafa, glob); `tumia crypto` (md5, sha256, kodeBase64, …); `tumia os` (toka, kimbiza); `tumia mfumo` (safishaMemori, takwimuMemori, …).

---

## Curriculum map (at a glance)

| Part | Theme | Units |
|------|--------|--------|
| 0 | Environment and first program | 0.1 |
| 1 | Values, types, and expressions | 1.1–1.5 |
| 2 | Control flow | 2.1–2.3 |
| 3 | Abstraction | 3.1–3.3 |
| 4 | Data structures | 4.1–4.5 |
| 5 | I/O and data formats | 5.1–5.3 |
| 6 | Text and time | 6.1–6.2 |
| 7 | Numerics | 7.1–7.2 |
| 8 | Networks and the web | 8.1–8.3 |
| 9 | System and environment | 9.1–9.4 |

**CS1 core** (typical first course): Parts 0–4 (through data structures and iteration).  
**CS1+** (same course with libraries): add Part 5 (files, JSON) and Part 6 (regex, time) as needed.  
**CS2-style / applied**: Parts 7–9 (math, big integers, HTTP, URLs, paths, crypto, OS, runtime).

---

## Reference and maintenance

- **Reference docs:** `repl/docs/en/` (26 .md files). REPL TOC in `repl/repl.go` omits range, sets, hisabati; adding them would align the in-REPL index with this curriculum.
- **Examples:** 31 scripts in `examples/`; table above maps them to units.
- **Policy:** `docs/POLICY.md` defines builtins, modules, and methods; units without a dedicated reference doc (url, http server, njia, crypto, mfumo, os) are covered there and/or by examples.
- **How-to pages:** `docs/howto/` holds learning units (e.g. 00-getting-started); numbering can be aligned to this part/unit scheme (e.g. 01-variables-and-types → 1.1).
