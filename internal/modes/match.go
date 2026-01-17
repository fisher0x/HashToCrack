package modes

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fisher0x/hashtocrack/internal/ntds"
	"github.com/fisher0x/hashtocrack/internal/utils"
)

// RunMatch matches NTDS entries with cracked passwords from a potfile
// Output format: username:hash:password:status
func RunMatch(ntdsFile, crackFile, outfile string, includeDisabled, includeMachines bool) {
	// Load potfile
	potfile, err := ntds.LoadPotfile(crackFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading potfile: %v\n", err)
		os.Exit(1)
	}

	// Open NTDS file
	file, err := os.Open(ntdsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening NTDS file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var output *os.File
	if outfile != "" {
		if err := utils.EnsureDir(outfile); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
		output, err = os.Create(outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := ntds.ParseLine(scanner.Text())
		if err != nil {
			continue
		}

		// Skip disabled unless flag is set
		if entry.IsDisabled && !includeDisabled {
			continue
		}

		// Skip machine accounts unless flag is set
		if entry.IsMachine && !includeMachines {
			continue
		}

		status := "Enabled"
		if entry.IsDisabled {
			status = "Disabled"
		}

		// Check if hash is cracked
		password := ""
		if pwd, found := potfile[strings.ToLower(entry.NTHash)]; found {
			password = pwd
		}

		var line string
		if password != "" {
			line = fmt.Sprintf("%s:%s:%s:%s", entry.Username, entry.NTHash, password, status)
		} else {
			line = fmt.Sprintf("%s:%s::%s", entry.Username, entry.NTHash, status)
		}

		if output != nil {
			fmt.Fprintln(output, line)
		} else {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if outfile != "" {
		fmt.Fprintf(os.Stderr, "[+] Matched results written to: %s\n", outfile)
	}
}
