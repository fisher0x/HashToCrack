package ntds

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fisher0x/hashtocrack/internal/utils"
)

// ParseLine parses a single line from NTDS file
func ParseLine(line string) (*Entry, error) {
	line = utils.CleanLine(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	entry := &Entry{RawLine: line}

	// Check if disabled - look for status=Disabled pattern
	lowerLine := strings.ToLower(line)
	entry.IsDisabled = strings.Contains(lowerLine, "status=disabled") ||
		(strings.Contains(lowerLine, "disabled") && !strings.Contains(lowerLine, "status=enabled"))

	// Remove status suffix for parsing - handle various formats
	// Format: ::: (status=Enabled) or ::: (status=Disabled)
	cleanedLine := line

	// Try to find and remove the status part
	if idx := strings.Index(line, "(status="); idx != -1 {
		cleanedLine = strings.TrimSpace(line[:idx])
	}

	// Remove trailing colons
	cleanedLine = strings.TrimRight(cleanedLine, ":")

	// Split by colon
	parts := strings.Split(cleanedLine, ":")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid format: expected at least 4 colon-separated fields, got %d", len(parts))
	}

	entry.Username = parts[0]
	entry.RID = parts[1]
	entry.LMHash = parts[2]
	entry.NTHash = parts[3]

	// Check if machine account (ends with $)
	entry.IsMachine = strings.HasSuffix(entry.Username, "$")

	return entry, nil
}

// ParseAnalyticsLine parses a line from analytics file (output from match mode)
func ParseAnalyticsLine(line string) (*CrackedEntry, error) {
	line = utils.CleanLine(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	parts := strings.Split(line, ":")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid format")
	}

	entry := &CrackedEntry{}
	entry.Username = parts[0]
	entry.NTHash = parts[1]

	// Password is everything between hash and status
	// Status is the last part
	status := parts[len(parts)-1]
	entry.IsDisabled = status == "Disabled"

	// Password is parts[2] to parts[len-2] joined (password might contain colons)
	if len(parts) > 4 {
		entry.Password = strings.Join(parts[2:len(parts)-1], ":")
	} else {
		entry.Password = parts[2]
	}

	entry.Cracked = entry.Password != ""
	entry.IsMachine = strings.HasSuffix(entry.Username, "$")

	return entry, nil
}

// IsAnalyticsFile checks if file is an analytics file (output from match mode)
func IsAnalyticsFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := utils.CleanLine(scanner.Text())
		// Analytics file format: username:hash:password:status
		parts := strings.Split(line, ":")
		if len(parts) >= 4 {
			lastPart := parts[len(parts)-1]
			return lastPart == "Enabled" || lastPart == "Disabled"
		}
	}
	return false
}
