// Command grist is the CLI entry point for the Grist recipe toolchain.
//
// Subcommands (skeleton — full behaviour lands in later phases):
//
//	grist run   <recipe.yaml>   Apply a recipe to its declared source.
//	grist build <recipe.yaml>   Generate Go source from a recipe.
//	grist tui                   Launch the interactive editor (Phase 3).
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// version is overridden at release time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(os.Stderr, "grist:", err)
		}
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		usage(stdout)
		return nil
	}

	switch args[0] {
	case "run":
		return cmdRun(args[1:], stdout, stderr)
	case "build":
		return cmdBuild(args[1:], stdout, stderr)
	case "tui":
		return cmdTUI(args[1:], stdout, stderr)
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, "grist", version)
		return nil
	case "help", "--help", "-h":
		usage(stdout)
		return nil
	default:
		usage(stderr)
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `grist — TUI and recipe runtime for gseq-table

Usage:
  grist <command> [arguments]

Commands:
  run     Apply a recipe to its declared source
  build   Generate Go source from a recipe
  tui     Launch the interactive editor (Phase 3)
  version Print version information
  help    Show this message

Run "grist <command> -h" for per-command flags.
`)
}

func cmdRun(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: grist run <recipe.yaml>")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("exactly one recipe path is required")
	}
	return fmt.Errorf("not implemented yet: run %s", fs.Arg(0))
}

func cmdBuild(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: grist build <recipe.yaml>")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("exactly one recipe path is required")
	}
	return fmt.Errorf("not implemented yet: build %s", fs.Arg(0))
}

func cmdTUI(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: grist tui")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	return errors.New("not implemented yet: tui (Phase 3)")
}
