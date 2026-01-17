package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fisher0x/hashtocrack/internal/modes"
	"github.com/fisher0x/hashtocrack/internal/ntds"
)

// Options holds all parsed command-line flags
type Options struct {
	NTDSFile  string
	CrackFile string
	OutFile   string
	Disabled  bool
	Machines  bool
	PassPol   bool
	Report    bool
}

// ParseArgs parses command-line arguments and returns Options
func ParseArgs(args []string) *Options {
	opts := &Options{}

	if len(args) == 0 {
		return opts
	}

	opts.NTDSFile = args[0]

	i := 1
	for i < len(args) {
		arg := args[i]
		switch arg {
		case "-disabled", "--disabled":
			opts.Disabled = true
		case "-machines", "--machines":
			opts.Machines = true
		case "-passpol", "--passpol":
			opts.PassPol = true
		case "-report", "--report":
			opts.Report = true
		case "-o", "-outfile", "--outfile":
			if i+1 < len(args) {
				opts.OutFile = args[i+1]
				i++
			}
		default:
			// If not a flag, it might be a crackfile
			if !strings.HasPrefix(arg, "-") && opts.CrackFile == "" {
				opts.CrackFile = arg
			}
		}
		i++
	}

	return opts
}

// Run executes the appropriate mode based on parsed options
func Run(opts *Options) {
	if opts.NTDSFile == "" {
		PrintUsage()
		os.Exit(1)
	}

	// Determine mode
	if opts.CrackFile != "" {
		// Check if crackfile exists
		if _, err := os.Stat(opts.CrackFile); err == nil {
			// Mode 2: Match mode (ntdsfile + crackfile)
			modes.RunMatch(opts.NTDSFile, opts.CrackFile, opts.OutFile, opts.Disabled, opts.Machines)
		} else {
			fmt.Fprintf(os.Stderr, "Error: File '%s' not found\n", opts.CrackFile)
			os.Exit(1)
		}
	} else if opts.PassPol {
		// Mode 3: Analytics mode (detected by -passpol flag)
		modes.RunAnalytics(opts.NTDSFile, opts.OutFile, opts.Disabled, opts.Machines, opts.PassPol, opts.Report)
	} else {
		// Check file content to determine if it's analytics or extract mode
		if ntds.IsAnalyticsFile(opts.NTDSFile) {
			modes.RunAnalytics(opts.NTDSFile, opts.OutFile, opts.Disabled, opts.Machines, opts.PassPol, opts.Report)
		} else {
			// Mode 1: Extract hashes mode
			modes.RunExtract(opts.NTDSFile, opts.OutFile, opts.Disabled, opts.Machines)
		}
	}
}
