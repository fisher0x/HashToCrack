package ntds

import (
	"bufio"
	"os"
	"strings"

	"github.com/fisher0x/hashtocrack/internal/utils"
)

// LoadPotfile loads hashcat potfile into a map (hash -> password)
func LoadPotfile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	potfile := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := utils.CleanLine(scanner.Text())
		if line == "" {
			continue
		}

		// Split at first colon (password might contain colons)
		idx := strings.Index(line, ":")
		if idx == -1 {
			continue
		}

		hash := strings.ToLower(line[:idx])
		password := line[idx+1:]
		potfile[hash] = password
	}

	return potfile, scanner.Err()
}
