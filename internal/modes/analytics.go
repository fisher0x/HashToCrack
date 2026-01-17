package modes

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/fisher0x/hashtocrack/internal/ntds"
	"github.com/fisher0x/hashtocrack/internal/utils"
)

// isComplexPassword checks if password meets DOMAIN_PASSWORD_COMPLEX requirements
// Requirements: >= 8 chars, 3 of 4 categories (upper, lower, digit, special)
func isComplexPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	categories := 0
	if hasUpper {
		categories++
	}
	if hasLower {
		categories++
	}
	if hasDigit {
		categories++
	}
	if hasSpecial {
		categories++
	}

	return categories >= 3
}

// redactPassword returns a redacted version of the password
// Shows first 3 characters, replaces the rest with asterisks
func redactPassword(password string) string {
	if len(password) <= 3 {
		return password
	}
	return password[:3] + strings.Repeat("*", len(password)-3)
}

// RunAnalytics generates statistics from matched file
func RunAnalytics(analyticsFile, outfile string, includeDisabled, includeMachines, showPasspol, redactPasswords bool) {
	file, err := os.Open(analyticsFile)
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

	// Collect statistics
	totalAccounts := 0
	crackedAccounts := 0
	complexCount := 0
	lengthDist := make(map[int]int)
	passwordCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := ntds.ParseAnalyticsLine(scanner.Text())
		if err != nil {
			continue
		}

		// Apply filters
		if entry.IsDisabled && !includeDisabled {
			continue
		}
		if entry.IsMachine && !includeMachines {
			continue
		}

		totalAccounts++

		if entry.Cracked {
			crackedAccounts++
			lengthDist[len(entry.Password)]++
			passwordCounts[entry.Password]++

			if isComplexPassword(entry.Password) {
				complexCount++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Calculate percentages
	crackPct := 0.0
	if totalAccounts > 0 {
		crackPct = float64(crackedAccounts) / float64(totalAccounts) * 100
	}

	complexPct := 0.0
	if crackedAccounts > 0 {
		complexPct = float64(complexCount) / float64(crackedAccounts) * 100
	}

	// Get top 10 passwords
	type pwdCount struct {
		Password string
		Count    int
	}
	var passwords []pwdCount
	for pwd, count := range passwordCounts {
		passwords = append(passwords, pwdCount{pwd, count})
	}
	sort.Slice(passwords, func(i, j int) bool {
		return passwords[i].Count > passwords[j].Count
	})
	if len(passwords) > 10 {
		passwords = passwords[:10]
	}

	// Sort length distribution
	var lengths []int
	for length := range lengthDist {
		lengths = append(lengths, length)
	}
	sort.Ints(lengths)

	// Generate report
	writeFunc := func(format string, args ...interface{}) {
		line := fmt.Sprintf(format, args...)
		if output != nil {
			fmt.Fprint(output, line)
		} else {
			fmt.Print(line)
		}
	}

	writeFunc("\n")
	writeFunc("╔══════════════════════════════════════════════════════════════╗\n")
	writeFunc("║           HASHTOCRACK - Password Analytics Report           ║\n")
	writeFunc("╚══════════════════════════════════════════════════════════════╝\n\n")

	// Filters applied
	writeFunc("Filters Applied:\n")
	writeFunc("  • Disabled accounts: %s\n", utils.BoolToIncluded(includeDisabled))
	writeFunc("  • Machine accounts:  %s\n", utils.BoolToIncluded(includeMachines))
	writeFunc("\n")

	// General Statistics
	writeFunc("═══════════════════════════════════════════════════════════════\n")
	writeFunc("                      GENERAL STATISTICS                        \n")
	writeFunc("═══════════════════════════════════════════════════════════════\n\n")
	writeFunc("  Total Accounts Analyzed:  %d\n", totalAccounts)
	writeFunc("  Passwords Cracked:        %d (%.2f%%)\n", crackedAccounts, crackPct)
	writeFunc("  Passwords Not Cracked:    %d (%.2f%%)\n", totalAccounts-crackedAccounts, 100-crackPct)
	writeFunc("\n")

	// Progress bar
	barWidth := 40
	filled := int(crackPct / 100 * float64(barWidth))
	writeFunc("  Crack Progress: [")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			writeFunc("█")
		} else {
			writeFunc("░")
		}
	}
	writeFunc("] %.1f%%\n\n", crackPct)

	// Password Length Distribution
	writeFunc("═══════════════════════════════════════════════════════════════\n")
	writeFunc("                  PASSWORD LENGTH DISTRIBUTION                  \n")
	writeFunc("═══════════════════════════════════════════════════════════════\n\n")

	if len(lengths) > 0 {
		maxCount := 0
		for _, count := range lengthDist {
			if count > maxCount {
				maxCount = count
			}
		}

		for _, length := range lengths {
			count := lengthDist[length]
			pct := float64(count) / float64(crackedAccounts) * 100
			barLen := int(float64(count) / float64(maxCount) * 30)
			bar := strings.Repeat("▓", barLen)
			writeFunc("  %2d chars: %-30s %4d (%5.1f%%)\n", length, bar, count, pct)
		}
	} else {
		writeFunc("  No cracked passwords to analyze.\n")
	}
	writeFunc("\n")

	// Top 10 Passwords
	writeFunc("═══════════════════════════════════════════════════════════════\n")
	writeFunc("                    TOP 10 MOST USED PASSWORDS                  \n")
	writeFunc("═══════════════════════════════════════════════════════════════\n\n")

	if len(passwords) > 0 {
		writeFunc("  %-4s  %-30s  %s\n", "Rank", "Password", "Count")
		writeFunc("  ────  ──────────────────────────────  ─────\n")
		for i, pwd := range passwords {
			displayPwd := pwd.Password
			if redactPasswords {
				displayPwd = redactPassword(displayPwd)
			}
			if len(displayPwd) > 28 {
				displayPwd = displayPwd[:25] + "..."
			}
			writeFunc("  #%-3d  %-30s  %d\n", i+1, displayPwd, pwd.Count)
		}
	} else {
		writeFunc("  No cracked passwords to analyze.\n")
	}
	writeFunc("\n")

	// Password Policy Compliance
	if showPasspol {
		writeFunc("═══════════════════════════════════════════════════════════════\n")
		writeFunc("              PASSWORD POLICY COMPLIANCE ANALYSIS               \n")
		writeFunc("═══════════════════════════════════════════════════════════════\n\n")

		writeFunc("  Policy: DOMAIN_PASSWORD_COMPLEX\n")
		writeFunc("  Requirements:\n")
		writeFunc("    • Minimum 8 characters\n")
		writeFunc("    • At least 3 of 4 categories:\n")
		writeFunc("      - Uppercase letters (A-Z)\n")
		writeFunc("      - Lowercase letters (a-z)\n")
		writeFunc("      - Digits (0-9)\n")
		writeFunc("      - Special characters (!@#$%%^&*...)\n\n")

		writeFunc("  Results:\n")
		writeFunc("    Compliant passwords:     %d (%.2f%% of cracked)\n", complexCount, complexPct)
		writeFunc("    Non-compliant passwords: %d (%.2f%% of cracked)\n", crackedAccounts-complexCount, 100-complexPct)
		writeFunc("\n")

		// Compliance bar
		writeFunc("  Compliance: [")
		compFilled := int(complexPct / 100 * float64(barWidth))
		for i := 0; i < barWidth; i++ {
			if i < compFilled {
				writeFunc("█")
			} else {
				writeFunc("░")
			}
		}
		writeFunc("] %.1f%%\n\n", complexPct)
	}

	writeFunc("═══════════════════════════════════════════════════════════════\n")
	writeFunc("                        END OF REPORT                           \n")
	writeFunc("═══════════════════════════════════════════════════════════════\n")

	if outfile != "" {
		fmt.Fprintf(os.Stderr, "[+] Analytics report written to: %s\n", outfile)
	}
}
