package server

import (
	"fmt"
	"os"
)

const usageStr = `
Description: Parse text counting words and sentence locations, this command can be evoked as
either a command line utility or as a stand alone server process.

Usage: clidemo [options...] [input_filename]

Server options:
    -N, --name NAME                  NAME of the server
    -p, --port PORT                  PORT to listen on (default: 49152)
    -n, --connections MAX            MAX server connections allowed (default: 1000)
    -X, --procs MAX                  MAX processor cores to use from the machine (default: 0)
	                                 Anything <= 0 is no change to the environment.
    -d, --debug                      Enable debugging output (default: false)

File input options:
    -f, --file FILE                  Process input FILE

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

    # Server mode activated on port 8080; 10 conns; 2 processors
    clidemo -N "San Francisco" -p 8080 -n 10 -X 2

	# File input using -f flag with debug option
	clidemo -f /tmp/inputfiles/foo/bar.txt -d > out.txt

	# Implicit file input (no -f flag needed)
	clidemo /tmp/inputfiles/foo/bar.txt > out.txt

	# Piping input
	cat /tmp/inputfiles/foo/bar.txt | clidemo > out.txt
`

// end help text

// PrintUsageAndExit is used to print out command line options.
func PrintUsageAndExit() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
