# Nuru VS Code Extension

Syntax highlighting, LSP (diagnostics, go-to-definition, hover, completion, outline), and run support for Nuru. Detects `.nr` and `.sw` files.

## Run current file

- **Command Palette** → **Nuru: Run current file** to run the active script with the Nuru interpreter.
- Ensure the Nuru binary is on your PATH, or set **Nuru: Interpreter Path** in settings.

## Language Server (LSP)

LSP is enabled by default (**Nuru: Enable Language Server**). To use it:

1. Build the LSP binary from the repo root: `go build -o nuru-lsp ./cmd/nuru-lsp`
2. Put `nuru-lsp` on your PATH, or set **Nuru: Language Server Path** in settings to the full path to the binary.
3. Open a `.nr` or `.sw` file; the extension will start the server when the language is activated.

You get diagnostics, go-to-definition, hover, completion, and document outline (symbols). If the server fails to start, check the path and ensure the binary runs; the extension will show an error message.

## Screenshot

Add `screenshot.png` under `extensions/vscode/assets/` for marketplace and README preview. The image is optional.

## Install (from repo)

The extension is not yet on the Marketplace. Install from the Nuru repo:

1. Clone [NuruProgramming/Nuru](https://github.com/NuruProgramming/Nuru) or download the [nuru extension folder](https://github.com/NuruProgramming/Nuru/tree/main/extensions/vscode/nuru).
2. Copy the **nuru** folder (`extensions/vscode/nuru`) into your VS Code extensions directory:
   - **Windows:** `%USERPROFILE%\.vscode\extensions`
   - **Linux / macOS:** `~/.vscode/extensions`
3. Run `npm install` inside the copied `nuru` folder (for `vscode-languageclient`).
4. Reload VS Code.

When the extension is published to the Marketplace, you will be able to install it from the Extensions view.
