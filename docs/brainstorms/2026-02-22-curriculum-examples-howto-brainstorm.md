# Brainstorm: Curriculum-Style Examples (Modular How-To)

**Date:** 2026-02-22  
**Topic:** Well-documented examples in the form of a computer-science-style curriculum, modular by topic.

---

## What We're Building

A **modular how-to curriculum** for Nuru: topic-based units that beginners, teachers, language learners, and self-learners can use in any order. Each unit is a single markdown file under `docs/howto/` containing learning outcomes, explained code examples (inline), and links to reference docs. No exercises; teachers and learners add their own tasks. **Every unit must have a runnable script** in `examples/`; each unit’s "Run this" section links to that script.

---

## Why This Approach

- **Audience:** Beginners + teachers, language learners + self-learners. Material must work for self-study and for instructors building a syllabus.
- **Modular by topic:** Units (e.g. Variables & types, Control flow, Files & JSON) can be combined or used standalone.
- **Examples only:** Fast to ship; no exercise/solution maintenance. Every unit has a runnable script in `examples/` so learners can run the code.
- **Location:** `docs/howto/` keeps curriculum identity clear and separate from reference (`repl/docs/`) and policy (`docs/POLICY.md`). Fits a "how to learn/teach" angle.

---

## Key Decisions

| Decision | Choice |
|----------|--------|
| Audience | Beginners + teachers, language learners + self-learners |
| Scope | Modular by topic (pick/combine units) |
| Content | Examples only (outcomes + explained examples + reference links) |
| Location | `docs/howto/` |
| Example code | Inline in markdown; each unit must link to a runnable `examples/<script>.nr` |
| Exercises | None in v1; teachers/learners add their own |

---

## Unit Template

Each unit file (e.g. `docs/howto/01-variables-and-types.md`) follows this structure:

1. **Title** (H1)
2. **Learning outcomes** — 2–4 bullet points (e.g. "Declare variables with `fanya`", "Use `andika` and `jaza` for I/O")
3. **Prerequisites** (optional) — e.g. "Install Nuru; run `nuru` for REPL or `nuru script.nr` for a file"
4. **Examples** — 2–5 blocks, each with:
   - Short heading (e.g. "Printing and reading input")
   - Code block (Nuru)
   - 1–3 sentence commentary (what it does, why it matters)
5. **Reference** — Links to relevant `repl/docs/en/` (and optionally `repl/docs/sw/`) pages
6. **Keywords** (optional) — Nuru/Swahili terms introduced (e.g. `andika`, `jaza`, `fanya`) for language learners
7. **Run this** (required) — Link to the unit’s runnable script in `examples/<name>.nr`; every unit must have a corresponding `.nr` file

---

## Suggested Initial Units

- `00-getting-started.md` — Install, REPL, run a file, first `andika("Habari")`
- `01-variables-and-types.md` — `fanya`, types, `aina`, `andika`/`jaza`
- `02-control-flow.md` — `kama`/`au kama`/`sivyo`, `kwa`/`wakati` loops
- `03-functions.md` — `unda`, parameters, `rudisha`
- `04-collections.md` — Arrays, dicts, sets, `mfululizo`, `seta`/`jozi`
- `05-files-and-json.md` — `tumia faili`, `tumia jsoni`, read/write and encode/decode
- `06-http-and-urls.md` — `tumia http` or `mtandao`, simple client/server or fetch (link to existing examples)

Numbering allows a suggested order; units remain usable in any order. More units (regex, time, crypto, etc.) can follow the same template.

---

## Index

`docs/howto/README.md` provides a short intro and a list of units (title + link + one-line description) so teachers and self-learners can see the full set and pick topics.

---

## Resolved Questions

- Primary audience: Beginners + teachers, language learners + self-learners
- Scope: Modular by topic
- Content: Examples only (no exercises in v1)
- Location: `docs/howto/` (Approach A)
- Example code: Inline in markdown; mandatory "Run this" link — every unit has a runnable script in `examples/`
- Language: English first; Kiswahili how-tos can be added later (e.g. `docs/howto/sw/`) when capacity allows
- Discoverability: README or CONTRIBUTING can link to "Learning path: [docs/howto/](docs/howto/)" for new contributors and learners
