package subcmd

import (
	"fmt"
	"io"
	"os"
)

// State describes the state of a subcommand.
// See the comments for Exitf to see how Interactive is used.
// It allows a program to run multiple commands.
type State struct {
	Name        string        // Name of the subcommand we are running.
	Interactive bool          // Whether the command is line-by-line.
	Stdin       io.Reader     // Where to read standard input
	Stdout      io.Writer     // Where to write standard output.
	Stderr      io.Writer     // Where to write error output.
	ExitCode    int           // Exit with non-zero status for minor problems.
}

// NewState returns a new State for the named subcommand.
func NewState(name string) *State {
	s := &State{Name: name}
	s.DefaultIO()
	return s
}

// Init initializes the config and client for the State.
func (s *State) Init() {

}

func (s *State) SetIO(stdin io.Reader, stdout, stderr io.Writer) {
	s.Stdin = stdin
	s.Stdout = stdout
	s.Stderr = stderr
}

func (s *State) DefaultIO() {
	s.SetIO(os.Stdin, os.Stdout, os.Stderr)
}

// Exitf prints the error and exits the program.
// If we are interactive, it calls panic("exit"), which is intended to be recovered
// from by the calling interpreter.
// We don't use log (although the packages we call do) because the errors
// are for regular people.
func (s *State) Exitf(format string, args ...interface{}) {
	format = fmt.Sprintf("gm: %s: %s\n", s.Name, format)
	fmt.Fprintf(s.Stderr, format, args...)
	if s.Interactive {
		panic("exit")
	}
	s.ExitCode = 1
	s.ExitNow()
}

// Exit calls s.Exitf with the error.
func (s *State) Exit(err error) {
	s.Exitf("%s", err)
}

// ExitNow terminates the process with the current ExitCode.
func (s *State) ExitNow() {
}
