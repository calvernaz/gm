package main

import (
	"flag"
)

func (s *State) config(args ...string) {
	const help = `
Config prints to standard output the contents of the current config file.
It works by saving the file at initialization time, so if the actual
file has changed since the command started, it will still show the
configuration being used.
`
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	
	// output stdout
	outFile := fs.String("out", "", "default standard output")
	// pretty
	pretty := fs.Bool("p", false, "show pretty output")
	
	s.ParseFlags(fs, args, help, "config")

	if *pretty {
		s.writePrettyOut(*outFile, s.configFile)
	} else {
		s.writeOut(*outFile, s.configFile)
	}
}
