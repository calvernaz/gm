package subcmd

import (
	"flag"
	"fmt"
	"os"
)

// ParseFlags parses the flags in the command line arguments,
// according to those set in the flag set.
func (s *State) ParseFlags(fs *flag.FlagSet, args []string, help, usage string) {
	helpFlag := fs.Bool("help", false, "print more information about the command")
	usageFn := func() {
		fmt.Fprintf(s.Stderr, "Usage: upspin %s\n", usage)
		if *helpFlag {
			fmt.Fprintln(s.Stderr, help)
		}
		// How many flags?
		n := 0
		fs.VisitAll(func(*flag.Flag) { n++ })
		if n > 0 {
			fmt.Fprintf(s.Stderr, "Flags:\n")
			fs.PrintDefaults()
		}
		if s.Interactive {
			panic("exit")
		}
	}
	fs.Usage = usageFn
	err := fs.Parse(args)
	if err != nil {
		s.Exit(err)
	}
	if *helpFlag {
		fs.Usage()
		os.Exit(2)
	}
}
