# Nuru VS Code Extension

Syntax highlighting, LSP (diagnostics, go-to-definition, hover, completion, outline, rename, code actions, formatting), debugger (DAP), and run support for Nuru. Detects `.nr` and `.sw` files.

## Run current file

- **Command Palette** → **Nuru: Run current file** to run the active script with the Nuru interpreter.
- Ensure the Nuru binary is on your PATH, or set **Nuru: Interpreter Path** in settings.

## Language Server (LSP)

LSP is enabled by default (**Nuru: Enable Language Server**). To use it:

1. Build the LSP binary from the repo root: `go build -o nuru-lsp ./cmd/nuru-lsp`
2. Put `nuru-lsp` on your PATH, or set **Nuru: Language Server Path** in settings to the full path to the binary.
3. Open a `.nr` or `.sw` file; the extension will start the server when the language is activated.

You get: diagnostics, go-to-definition, hover, completion, document outline (symbols), rename, code actions (quick fix, organize imports), and document formatting. If the server fails to start, check the path and ensure the binary runs; the extension will show an error message.

## Debugger (DAP)

To debug a Nuru script:

1. Build the DAP binary from the repo root: `go build -o nuru-dap ./cmd/nuru-dap`
2. Put `nuru-dap` on your PATH, or set **Nuru: Debug Adapter Path** in settings to the full path to the binary.
3. Open a `.nr` file, add breakpoints (click gutter or F9), then **Run → Start Debugging** (F5) or use the **Run and Debug** view with the **Launch Nuru file** configuration.

You get: breakpoints, call stack, local variables, and continue/step controls. The launch config runs the current file; set `program` to `${file}` (default) or another script path.

## Documentation

Full LSP and DAP documentation (capabilities, implementation notes, source layout): [docs/tooling/LSP_AND_DAP.md](../../docs/tooling/LSP_AND_DAP.md) in the Nuru repo.

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
