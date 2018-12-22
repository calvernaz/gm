package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	
	"github.com/calvernaz/gm/flags"
	"github.com/calvernaz/gm/gm"
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
	commands = map[string]func(*State, ...string){
		"config": (*State).config,
		"add":    (*State).add,
		"update": (*State).update,
	}
)

type State struct {
	*subcmd.State
	configFile []byte // The contents of the config file we loaded.

	gmc        *gm.GitManagerConfig
}

func main() {
	state, args, ok := setup(flag.CommandLine, os.Args[1:])
	if !ok || len(args) == 0 {
		help()
	}
	if args[0] == "help" {
		help(args[1:]...)
	}

	// loop accepting and running operations
	go state.gmc.Loop()
	
	// run the command
	state.run(args)
	
	// wait until operation  is done
	state.gmc.Wait()
	
	state.ExitNow()
}

// setup initializes the gm command given the full command-line argument
// list, args. It applies any global flags set on the command line and returns
// the initialized State and the arg list after the global flags, starting with
// the subcommand ("ls", "info", etc.) that will be run.
func setup(fs *flag.FlagSet, args []string) (*State, []string, bool) {
	log.SetFlags(0)
	log.SetPrefix("gm: ")
	fs.Usage = usage
	flags.ParseArgsInto(fs, args, flags.Client, "version")
	if flags.Version {
		_, _ = fmt.Fprintln(os.Stdout, version)
		os.Exit(2)
	}
	if len(fs.Args()) < 1 {
		return nil, nil, false
	}
	state := newState(strings.ToLower(fs.Arg(0)))
	state.gmc = &gm.GitManagerConfig{Version: version}
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
		_, _ = fmt.Fprintln(os.Stderr, intro)
	} else {
		// Simplest solution is re-execing.
		command := exec.Command("gm", cmd, "-help")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		_ = command.Run()
	}
	os.Exit(2)
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\n\tgm [globalflags] <command> [flags]\n")
	printCommands()
	_, _ = fmt.Fprintf(os.Stderr, "\nGlobal flags:\n")
	flag.PrintDefaults()
}

// printCommands shows the available commands
func printCommands() {
	_, _ = fmt.Fprintf(os.Stderr, "\nCommands:\n")
	var cmdStrs []string
	for cmd := range commands {
		cmdStrs = append(cmdStrs, cmd)
	}
	
	// Display "shell" first as it's not in "commands".
	_, _ = fmt.Fprintf(os.Stderr, "\tshell (Interactive mode)\n")
	sort.Strings(cmdStrs)
	// There may be dups; filter them.
	prev := ""
	for _, cmd := range cmdStrs {
		if cmd == prev {
			continue
		}
		prev = cmd
		_, _ = fmt.Fprintf(os.Stderr, "\t%s\n", cmd)
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
	// open the git manager config
	err := s.gmc.Open("")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		s.gmc.Close()
		return
	}
	
	data, err := ioutil.ReadAll(s.gmc.File())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\tError reading config file: %v", err)
		return
	}
	s.configFile = data
}

// getCommand looks up the command named by op.
func (s *State) getCommand(op string) func(*State, ...string) {
	op = strings.ToLower(op)
	fn := commands[op]
	if fn != nil {
		return fn
	}
	path, err := exec.LookPath("gm-" + op)
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
func (s *State) run(args []string) {
	cmd := s.getCommand(args[0])
	cmd(s, args[1:]...)
}

// writeOut writes to the named file or to stdout if it is empty
func (s *State) writeOut(file string, data []byte) {
	// Write to outfile or to stdout if none set
	if file == "" {
		_, err := s.Stdout.Write(data)
		if err != nil {
			s.Exitf("copying to output failed: %v", err)
		}
		return
	}
}

func (s *State) writePrettyOut(file string, data []byte) {
	if file == "" {
		var dat map[string]interface{}
		if err := json.Unmarshal(data, &dat); err != nil {
			s.Exitf("failed to output: %v", err)
		}
		
		b, err := json.MarshalIndent(dat, "", "  ")
		if err != nil {
			s.Exitf("failed to output: %v", err)
		}
		
		_, err = s.Stdout.Write(b)
		if err != nil {
			s.Exitf("copying to output failed: %v", err)
		}
		
		return
	}
}
// usageAndExit prints usage message from provided FlagSet,
// and exits the program with status code 2.
func usageAndExit(fs *flag.FlagSet) {
	fs.Usage()
	os.Exit(2)
}
