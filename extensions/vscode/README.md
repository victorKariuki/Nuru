# Nuru VSCode Extension

Syntax highlighting, LSP (diagnostics, go-to-definition, hover, completion), and run support for Nuru. Detects `.nr` and `.sw` files.

## Run current file

- **Command Palette** → **Nuru: Run current file** to run the active script with the Nuru interpreter.
- Ensure the Nuru binary is on your PATH, or set **Nuru: Interpreter Path** in settings.

## Language Server (LSP)

To enable diagnostics, go-to-definition, hover, and completion:

1. Build the LSP binary from the repo root: `go build -o nuru-lsp ./cmd/nuru-lsp`
2. Put `nuru-lsp` on your PATH, or set **Nuru: Language Server Path** in settings to the full path to the binary.
3. Open a `.nr` or `.sw` file; the extension will start the server when the language is activated.

## Screenshots

A screenshot can be placed in `extensions/vscode/assets/screenshot.png` for the README and marketplace.

<p align="center">
<img alt="Nuru Programming Language" src="assets/screenshot.png">
</p>

## How To Install

### Download From Market Place

- Simply download the Nuru Extension from VSCode Market Place

### Windows

- Copy the whole [nuru folder](https://github.com/NuruProgramming/Nuru/tree/main/extensions/vscode/nuru) and paste it in the VSCode extensions directory found in `%USERPROFILE%\.vscode\extensions`
- Restart VSCode

### Linux and MacOS

- Copy the whole [nuru folder](https://github.com/NuruProgramming/Nuru/tree/main/extensions/vscode/nuru) and paste it in the VSCode extensions directory found in `~/.vscode/extensions`
- Restart VSCode
