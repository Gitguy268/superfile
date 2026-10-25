package utils

import (
	"os"

	"github.com/charmbracelet/x/term"
)

// TerminalOutput returns the file superfile should draw its UI on. It is stdout,
// unless stdout is redirected (as in `cd "$(spf --print-last-dir)"`) while
// stderr is still a terminal. Then the UI goes to stderr, so stdout only
// carries what superfile prints on exit.
func TerminalOutput() *os.File {
	if !term.IsTerminal(os.Stdout.Fd()) && term.IsTerminal(os.Stderr.Fd()) {
		return os.Stderr
	}
	return os.Stdout
}
