package main

import (
	"os"

	"github.com/jacobmaizel/go-interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
