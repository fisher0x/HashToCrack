package modes

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fisher0x/hashtocrack/internal/ntds"
	"github.com/fisher0x/hashtocrack/internal/utils"
)

// RunExtract extracts hashes from NTDS file
// This is equivalent to: grep -iv disabled ntdsfile | cut -d ':' -f4
// Or with -disabled flag: cat ntdsfile | cut -d ':' -f4
func RunExtract(ntdsFile, outfile string, includeDisabled, includeMachines bool) {
	file, err := os.Open(ntdsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
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

		line := entry.NTHash

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
		fmt.Fprintf(os.Stderr, "[+] Hashes written to: %s\n", outfile)
	}
}
