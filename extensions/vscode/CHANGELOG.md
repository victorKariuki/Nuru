# Changelog

All notable changes to the Nuru VS Code extension are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.0.1] - Unreleased

### Added

- Syntax highlighting for `.nr` and `.sw` files
- Language configuration (comments, brackets, auto-closing pairs)
- Nuru Language Server (LSP) support: diagnostics, go-to-definition, hover, completion, document outline (symbols), rename, code actions, formatting
- **Nuru: Run current file** command (Command Palette) to run the active script with the Nuru interpreter
- Snippets: `fanya`, `kama` / `au kama` / `sivyo`, `kwa`, `wakati`, `unda`, `tumia`, `badili` (switch), `jaribu` / `jaribu-bila` (try-catch), `rudisha`, `pakeji`, `vunja`, `endelea`
- Setting `nuru.enableLanguageServer` to turn LSP on or off (default: true)
- Setting `nuru.languageServerPath` for the LSP binary
- Setting `nuru.interpreterPath` for running scripts
- Nuru Debug Adapter (DAP): launch, breakpoints, call stack, scopes, variables, continue/step; setting `nuru.debugAdapterPath` for the DAP binary
- Client-side error message when the language server fails to start; run command warns if interpreter path is empty
- Documentation: [docs/tooling/LSP_AND_DAP.md](../../docs/tooling/LSP_AND_DAP.md) for full LSP and DAP reference

### Changed

- Grammar: control keywords, builtins, and constants aligned with language (token.go, POLICY)
- Grammar: operators `!=` and `**` added
- Grammar: builtins highlighted separately from function declarations
- Language configuration: valid JSON only (no JavaScript-style comments)
