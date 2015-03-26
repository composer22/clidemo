// clidemo is a simple demonstration of reading CLI info and evoking piping, file, or server work. This application is also used to POC features of Golang as needed.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/composer22/clidemo/parser"
	"github.com/composer22/clidemo/server"
)

// configureLogging configures logging for the application run.
func configureLogging() {
	log.SetFlags(log.Lmicroseconds)
}

// configureServerEnvironment configures the physical and logical server components for the application run.
func configureServerEnvironment(opts *server.Options) {
	if opts.MaxProcs > 0 {
		runtime.GOMAXPROCS(opts.MaxProcs)
	}
	log.Printf("INFO: NumCPU %d GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(-1))
}

// main is the main entry point for the application or server launch.
func main() {
	opts := server.Options{}
	var showVersion bool
	var fileIn string

	flag.StringVar(&opts.Name, "N", "", "Name of the server")
	flag.StringVar(&opts.Name, "--name", "", "Name of the server")
	flag.IntVar(&opts.Port, "P", server.DefaultPort, "Port to listen on (default: 49152)")
	flag.IntVar(&opts.Port, "--port", server.DefaultPort, "Port to listen on (default: 49152)")
	flag.IntVar(&opts.MaxConn, "n", server.DefaultMaxConnections,
		"Maximum server connections allowed (default: 4)")
	flag.IntVar(&opts.MaxConn, "--connections", server.DefaultMaxConnections,
		"Maximum server connections allowed (default: 4)")
	flag.IntVar(&opts.MaxProcs, "X", server.DefaultMaxProcs,
		"Maximum processor cores to use from the machine (default: <= 0 is no change")
	flag.IntVar(&opts.MaxProcs, "--procs", server.DefaultMaxProcs,
		"Maximum processor cores to use from the machine (default: <= 0 is no change)")
	flag.BoolVar(&opts.Debug, "d", false, "Enable debugging output (default: false)")
	flag.BoolVar(&opts.Debug, "--debug", false, "Enable debugging output (default: false)")
	flag.StringVar(&fileIn, "f", "", "Process input file")
	flag.StringVar(&fileIn, "--file", "", "Process input file")
	flag.BoolVar(&showVersion, "V", false, "Show version")
	flag.BoolVar(&showVersion, "--version", false, "Show version")
	flag.Usage = server.PrintUsageAndExit
	flag.Parse()

	// Version flag request?
	if showVersion {
		server.PrintVersionAndExit()
	}

	// Check additional params beyond the flags, such as commands or filename w/o -f.
	for _, arg := range flag.Args() {
		switch strings.ToLower(arg) {
		case "version":
			server.PrintVersionAndExit()
		case "help":
			server.PrintUsageAndExit()
		default: // input filename via w/o -f flag e.g. appname /tmp/foo/bar.txt.
			fileIn = arg
		}
	}

	// Get any stats we need for checking piped input.
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalln("[FATAL] os.Stdin.Stat(): ", err)
	}

	configureLogging()

	// Lets do work as a service or on direct input.
	switch {

	// Piped input text (higher priority than command line file names or server mode).
	case fi.Mode()&os.ModeNamedPipe == 1:
		p := parser.New()
		p.Execute(bufio.NewReader(os.Stdin))
		fmt.Print(p)

	// File input text higher priority than server mode.
	case fileIn != "":
		fi, err := os.Open(fileIn)
		if err != nil {
			log.Fatalln("[FATAL] Cannot open file ", fileIn, ": ", err)
		}
		defer fi.Close()
		p := parser.New()
		p.Execute(bufio.NewReader(fi))
		fmt.Print(p)

	// Server mode.
	default:
		configureServerEnvironment(&opts)
		s := server.New(&opts)
		s.Start()
	}
}
