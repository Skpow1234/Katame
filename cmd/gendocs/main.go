// Command gendocs generates man pages for katame from the cobra command tree.
package main

import (
	"fmt"
	"os"

	"github.com/Skpow1234/Katame/internal/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	out := "man"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", out, err)
		os.Exit(1)
	}
	header := &doc.GenManHeader{
		Title:   "KATAME",
		Section: "1",
		Source:  "Katame",
	}
	if err := doc.GenManTree(cli.NewRootCmd(), header, out); err != nil {
		fmt.Fprintf(os.Stderr, "gen man: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("wrote man pages to %s/\n", out)
}
