package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/calvernaz/gm"
	"github.com/calvernaz/gm/flags"
	"github.com/calvernaz/gm/subcmd"
)

const (
	version = "1"
	intro   = `
For a list of available subcommands and global flags, run

        gm -help
`
)

var (
	configPath = filepath.Join("gm", "gm.json")
	commands   = map[string]func(*State, ...string){}
)

type State struct {
	*subcmd.State
	configFile []byte // The contents of the config file we loaded.
}

func main() {
	gmc := gm.GitManagerConfig{Version: version}

	// open the git manager config
	err := gmc.Open(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		gmc.Close()
		return
	}

	state, args, ok := setup(flag.CommandLine, os.Args[1:])
	if !ok || len(args) == 0 {
		help()
	}
	if args[0] == "help" {
		help(args[1:]...)
	}

	state.run(args)
	state.ExitNow()
}

// setup initializes the upspin command given the full command-line argument
// list, args. It applies any global flags set on the command line and returns
// the initialized State and the arg list after the global flags, starting with
// the subcommand ("ls", "info", etc.) that will be run.
func setup(fs *flag.FlagSet, args []string) (*State, []string, bool) {
	log.SetFlags(0)
	log.SetPrefix("gm: ")
	fs.Usage = usage
	flags.ParseArgsInto(fs, args, flags.Client, "version")
	if flags.Version {
		fmt.Fprint(os.Stdout, version)
		os.Exit(2)
	}
	if len(fs.Args()) < 1 {
		return nil, nil, false
	}
	state := newState(strings.ToLower(fs.Arg(0)))
	state.init()

	return state, fs.Args(), true
}

// help prints the help for the arguments provided, or if there is none,
// for the command itself.
func help(args ...string) {
	// Find the first non-flag argument.
	cmd := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			cmd = arg
			break
		}
	}
	if cmd == "" {
		fmt.Fprintln(os.Stderr, intro)
	} else {
		// Simplest solution is re-execing.
		command := exec.Command("gm", cmd, "-help")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Run()
	}
	os.Exit(2)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of gm:\n")
	fmt.Fprintf(os.Stderr, "\tgm [globalflags] <command> [flags]\n")
	printCommands()
	fmt.Fprintf(os.Stderr, "Global flags:\n")
	flag.PrintDefaults()
}

// printCommands shows the available commands, including those installed
// as separate binaries called "upspin-foo".
func printCommands() {
	fmt.Fprintf(os.Stderr, "gm commands:\n")
	var cmdStrs []string
	for cmd := range commands {
		cmdStrs = append(cmdStrs, cmd)
	}

	// Display "shell" first as it's not in "commands".
	fmt.Fprintf(os.Stderr, "\tshell (Interactive mode)\n")
	sort.Strings(cmdStrs)
	// There may be dups; filter them.
	prev := ""
	for _, cmd := range cmdStrs {
		if cmd == prev {
			continue
		}
		prev = cmd
		fmt.Fprintf(os.Stderr, "\t%s\n", cmd)
	}
}

// newState returns a State with enough initialized to run exit, etc.
// It does not contain a Config.
func newState(name string) *State {
	s := &State{
		State: subcmd.NewState(name),
	}
	return s
}

// init initializes the State with what is required to run the subcommand,
// usually including setting up a Config.
func (s *State) init() {

}

// getCommand looks up the command named by op.
// If it's in the commands tables, we're done.
// If not, it looks for a binary with the equivalent name
// (upspin foo is implemented by upspin-foo).
// If the command still can't be found, it exits after listing the
// commands that do exist.
func (s *State) getCommand(op string) func(*State, ...string) {
	op = strings.ToLower(op)
	fn := commands[op]
	if fn != nil {
		return fn
	}
	path, err := exec.LookPath("upspin-" + op)
	if err == nil {
		return func(s *State, args ...string) {
			s.runCommand(path, append(flags.Args(), args...)...)
		}
	}
	printCommands()
	s.Exitf("no such command %q", op)
	return nil
}

func (s *State) runCommand(path string, args ...string) {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		s.Exit(err)
	}
}

// run runs a single command specified by the arguments, which should begin with
// the subcommand ("ls", "info", etc.).
func (state *State) run(args []string) {
	cmd := state.getCommand(args[0])
	cmd(state, args[1:]...)
}
