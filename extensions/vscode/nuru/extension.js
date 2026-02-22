"use strict";

const path = require("path");
const vscode = require("vscode");
const {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} = require("vscode-languageclient/node");

/** @type {LanguageClient | undefined} */
let client;

function activate(context) {
  const config = vscode.workspace.getConfiguration("nuru");
  const enableLsp = config.get("enableLanguageServer", true);

  context.subscriptions.push(
    vscode.debug.registerDebugAdapterDescriptorFactory("nuru", {
      createDebugAdapterDescriptor() {
        const cfg = vscode.workspace.getConfiguration("nuru");
        const adapterPath = cfg.get("debugAdapterPath", "nuru-dap");
        return new vscode.DebugAdapterExecutable(adapterPath, []);
      },
    })
  );

  if (enableLsp) {
    const serverCommand = config.get("languageServerPath", "nuru-lsp");
    const serverOptions = { command: serverCommand, args: [] };
    const clientOptions = {
      documentSelector: [{ scheme: "file", language: "nuru" }],
    };
    client = new LanguageClient(
      "nuruLsp",
      "Nuru Language Server",
      serverOptions,
      clientOptions
    );
    client.start().then(() => {
      context.subscriptions.push(client);
    }).catch((err) => {
      const msg = err && err.message ? err.message : String(err);
      vscode.window.showErrorMessage(
        "Nuru: Language Server failed to start. " +
        "Ensure nuru-lsp is on your PATH or set \"Nuru: Language Server Path\" in settings. " +
        msg
      );
    });
  }

  context.subscriptions.push(
    vscode.commands.registerCommand("nuru.runFile", async () => {
      const editor = vscode.window.activeTextEditor;
      if (!editor || editor.document.languageId !== "nuru") {
        vscode.window.showWarningMessage(
          "Open a Nuru (.nr or .sw) file to run."
        );
        return;
      }
      const file = editor.document.uri.fsPath;
      const interpreter = config.get("interpreterPath", "nuru");
      if (!interpreter || String(interpreter).trim() === "") {
        vscode.window.showWarningMessage(
          "Nuru: Interpreter path is empty. Set \"Nuru: Interpreter Path\" in settings or put the nuru binary on your PATH."
        );
        return;
      }
      const term = vscode.window.createTerminal({
        name: "Nuru",
        cwd: path.dirname(file),
      });
      term.show();
      term.sendText(`${interpreter} "${file}"`);
    })
  );
}

function deactivate() {
  if (client) {
    return client.stop();
  }
}

module.exports = { activate, deactivate };
