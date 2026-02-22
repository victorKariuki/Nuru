"use strict";

const vscode = require("vscode");
const {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} = require("vscode-languageclient/node");

/** @type {LanguageClient | undefined} */
let client;

function activate(context) {
  const serverCommand = vscode.workspace
    .getConfiguration("nuru")
    .get("languageServerPath", "nuru-lsp");

  const serverOptions = {
    command: serverCommand,
    args: [],
  };

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
  });

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
      const interpreter = vscode.workspace
        .getConfiguration("nuru")
        .get("interpreterPath", "nuru");
      const term = vscode.window.createTerminal({
        name: "Nuru",
        cwd: require("path").dirname(file),
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
