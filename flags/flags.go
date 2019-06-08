package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/calvernaz/gm/manager"
)

const (
	defaultLog        = "info"
)

// flagVar represents a flag in this package.
type flagVar struct {
	set  func(fs *flag.FlagSet) // Set the value at parse time.
	arg  func() string          // Return the argument to set the flag.
	arg2 func() string          // Return the argument to set the second flag; usually nil.
}

// None is the set of no flags. It is rarely needed as most programs
// use either the Server or Client set.
var None []string

// Server is the set of flags most useful in servers. It can be passed as the
// argument to Parse to set up the package for a server.
var Server = []string{
	"config", "log", "http", "https", "letscache", "tls", "addr", "insecure",
}

// Client is the set of flags most useful in clients. It can be passed as the
// argument to Parse to set up the package for a client.
var Client = []string{
	"log",
}

var (
	// Log ("log") sets the level of logging (implements flag.Value).
	Log logFlag

	// Version causes the program to print its release version and exit.
	// The printed version is only meaningful in released binaries.
	Version = false
)


// flags is a map of flag registration functions keyed by flag name,
// used by Parse to register specific (or all) flags.
var flags = map[string]*flagVar{
	"version": {
		set: func(fs *flag.FlagSet) {
			fs.BoolVar(&Version, "version", false, "print build version and exit")
		},
		arg: func() string {
			if !Version {
				return ""
			}
			return "-version"
		},
	},
	"log": &flagVar{
		set: func(fs *flag.FlagSet) {
			Log.Set("info")
			fs.Var(&Log, "log", "`level` of logging: debug, info, error, disabled")
		},
		arg: func() string { return strArg("log", Log.String(), defaultLog) },
	},
}

type logFlag string

// String implements flag.Value.
func (f logFlag) String() string {
	return string(f)
}

// Set implements flag.Value.
func (f *logFlag) Set(level string) error {
	err := manager.SetLevel(level)
	if err != nil {
		return err
	}
	*f = logFlag(manager.GetLevel())
	return nil
}

// Get implements flag.Getter.
func (logFlag) Get() interface{} {
	return manager.GetLevel()
}

// Parse registers the command-line flags for the given default flags list, plus
// any extra flag names, and calls flag.Parse. Passing no flag names in either
// list registers all flags. Passing an unknown name triggers a panic.
// The Server and Client variables contain useful default sets.
//
// Examples:
// 	flags.Parse(flags.Client) // Register all client flags.
//	flags.Parse(flags.Server, "cachedir") // Register all server flags plus cachedir.
// 	flags.Parse(nil) // Register all flags.
// 	flags.Parse(flags.None, "config", "endpoint") // Register only config and endpoint.
func Parse(defaultList []string, extras ...string) {
	ParseArgsInto(flag.CommandLine, os.Args[1:], defaultList, extras...)
}

// ParseInto is the same as Parse but accepts a FlagSet argument instead of
// using the default flag.CommandLine FlagSet.
func ParseInto(fs *flag.FlagSet, defaultList []string, extras ...string) {
	ParseArgsInto(fs, os.Args[1:], defaultList, extras...)
}

// ParseArgs is the same as Parse but uses the provided argument list
// instead of those provided on the command line. For ParseArgs, the
// initial command name should not be provided.
func ParseArgs(args, defaultList []string, extras ...string) {
	ParseArgsInto(flag.CommandLine, args, defaultList, extras...)
}

// ParseArgsInto is the same as ParseArgs but accepts a FlagSet argument instead of
// using the default flag.CommandLine FlagSet.
func ParseArgsInto(fs *flag.FlagSet, args, defaultList []string, extras ...string) {
	if len(defaultList) == 0 && len(extras) == 0 {
		RegisterInto(fs)
	} else {
		if len(defaultList) > 0 {
			RegisterInto(fs, defaultList...)
		}
		if len(extras) > 0 {
			RegisterInto(fs, extras...)
		}
	}
	fs.Parse(args)
}

// Register registers the command-line flags for the given flag names.
// Unlike Parse, it may be called multiple times.
// Passing zero names install all flags.
// Passing an unknown name triggers a panic.
//
// For example:
// 	flags.Register("config", "endpoint") // Register Config and Endpoint.
// or
// 	flags.Register() // Register all flags.
func Register(names ...string) {
	RegisterInto(flag.CommandLine, names...)
}

// RegisterInto  is the same as Register but accepts a FlagSet argument instead of
// using the default flag.CommandLine FlagSet.
func RegisterInto(fs *flag.FlagSet, names ...string) {
	if len(names) == 0 {
		// Register all flags if no names provided.
		for _, f := range flags {
			f.set(fs)
		}
	} else {
		for _, n := range names {
			f, ok := flags[n]
			if !ok {
				panic(fmt.Sprintf("unknown flag %q", n))
			}
			f.set(fs)
		}
	}
}

// Args returns a slice of -flag=value strings that will recreate
// the state of the flags. Flags set to their default value are elided.
func Args() []string {
	var args []string
	for _, f := range flags {
		arg := f.arg()
		if arg == "" {
			continue
		}
		args = append(args, arg)
		if f.arg2 != nil {
			args = append(args, f.arg2())
		}
	}
	return args
}

// strVar returns a flagVar for the given string flag.
func strVar(value *string, name, _default, usage string) *flagVar {
	return &flagVar{
		set: func(fs *flag.FlagSet) {
			fs.StringVar(value, name, _default, usage)
		},
		arg: func() string {
			return strArg(name, *value, _default)
		},
	}
}

// strArg returns a command-line argument that will recreate the flag,
// or the empty string if the value is the default.
func strArg(name, value, _default string) string {
	if value == _default {
		return ""
	}
	return "-" + name + "=" + value
}

// String implements flag.Value.

// String implements flag.Value.

// Set implements flag.Value.

// Get implements flag.Getter.
