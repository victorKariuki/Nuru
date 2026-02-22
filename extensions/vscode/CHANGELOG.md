# Changelog

All notable changes to the Nuru VS Code extension are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.0.1] - Unreleased

### Added

- Syntax highlighting for `.nr` and `.sw` files
- Language configuration (comments, brackets, auto-closing pairs)
- Nuru Language Server (LSP) support: diagnostics, go-to-definition, hover, completion
- **Nuru: Run current file** command (Command Palette) to run the active script with the Nuru interpreter
- Snippets: `fanya`, `kama` / `au kama` / `sivyo`, `kwa`, `wakati`, `unda`, `tumia`, `badili` (switch), `jaribu` / `jaribu-bila` (try-catch), `rudisha`, `pakeji`, `vunja`, `endelea`
- Setting `nuru.languageServerPath` for the LSP binary
- Setting `nuru.interpreterPath` for running scripts

### Changed

- Grammar: control keywords, builtins, and constants aligned with language (token.go, POLICY)
- Grammar: operators `!=` and `**` added
- Grammar: builtins highlighted separately from function declarations
- Language configuration: valid JSON only (no JavaScript-style comments)
