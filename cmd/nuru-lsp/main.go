package main

import (
	"os"

	"github.com/NuruProgramming/Nuru/lsp"
)

func main() {
	server := lsp.NewServer()
	if err := server.Run(os.Stdin, os.Stdout); err != nil {
		os.Exit(1)
	}
}
