# HashToCrack 🔐

A cross-platform NTDS hash analyzer and password statistics tool for security professionals.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com)
[![Release](https://img.shields.io/github/v/release/fisher0x/HashToCrack?include_prereleases)](https://github.com/fisher0x/HashToCrack/releases)

## Overview

HashToCrack helps security professionals analyze NTDS (Active Directory) parsed files and hashcat potfiles. It can:

- 🔑 **Extract** NT hashes from NTDS dumps for cracking
- 🔗 **Match** cracked hashes with their account owners
- 📊 **Analyze** password statistics and policy compliance
- 📝 **Report** with redacted passwords for safe sharing

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/fisher0x/HashToCrack/releases) page.

Available binaries:
- `HashToCrack_linux_amd64` - Linux 64-bit
- `HashToCrack_linux_arm64` - Linux ARM64
- `HashToCrack_darwin_amd64` - macOS Intel
- `HashToCrack_darwin_arm64` - macOS Apple Silicon
- `HashToCrack_windows_amd64.exe` - Windows 64-bit

### From Source

```bash
# Clone the repository
git clone https://github.com/fisher0x/HashToCrack.git
cd HashToCrack

# Build for your platform
go build -o HashToCrack ./cmd/HashToCrack

# Or build for all platforms
make all
```

### Using Go Install

```bash
go install github.com/fisher0x/hashtocrack/cmd/HashToCrack@latest
```

## Quick Start

```bash
# Extract hashes for hashcat
HashToCrack NTDS.dit -o hashes.txt

# Crack with hashcat
hashcat -m 1000 hashes.txt wordlist.txt -o potfile.txt

# Match passwords with accounts
HashToCrack NTDS.dit potfile.txt -o matched.txt

# Generate analytics report
HashToCrack matched.txt -passpol -report -o report.txt
```

## Usage

HashToCrack operates in three modes:

### 1. Extract Mode - Extract Hashes

Extract NT hashes from NTDS file for cracking with hashcat:

```bash
HashToCrack <ntdsfile> [-disabled] [-machines] [-o <outfile>]
```

| Flag | Description |
|------|-------------|
| `-disabled` | Include disabled accounts |
| `-machines` | Include machine accounts (ending with `$`) |
| `-o` | Write output to file |

**Examples:**
```bash
HashToCrack NTDS.dit                      # Extract enabled user hashes
HashToCrack NTDS.dit -disabled            # Include disabled accounts
HashToCrack NTDS.dit -machines            # Include machine accounts
HashToCrack NTDS.dit -disabled -machines -o hashes.txt
```

### 2. Match Mode - Match Hashes with Passwords

Match NTDS entries with a hashcat potfile to identify cracked accounts:

```bash
HashToCrack <ntdsfile> <potfile> [-disabled] [-machines] [-o <outfile>]
```

**Examples:**
```bash
HashToCrack NTDS.dit potfile.txt                    # Match and display
HashToCrack NTDS.dit potfile.txt -disabled          # Include disabled accounts
HashToCrack NTDS.dit potfile.txt -o matched.txt     # Save to file
```

**Output format:** `username:hash:password:status`

```
DOMAIN\jsmith:b4b9b02e6f09a9bd760f388b67351e2b:Summer2024!:Enabled
DOMAIN\admin:aad3b435b51404eeaad3b435b51404ee::Enabled
DOMAIN\olduser:5f4dcc3b5aa765d61d8327deb882cf99:password123:Disabled
```

### 3. Analytics Mode - Generate Statistics

Analyze matched results to generate comprehensive password statistics:

```bash
HashToCrack <matchedfile> [-disabled] [-machines] [-passpol] [-report] [-o <outfile>]
```

| Flag | Description |
|------|-------------|
| `-disabled` | Include disabled accounts in statistics |
| `-machines` | Include machine accounts in statistics |
| `-passpol` | Show password policy compliance analysis |
| `-report` | Redact passwords (show first 3 chars only) |
| `-o` | Write report to file |

**Examples:**
```bash
HashToCrack matched.txt                              # Basic statistics
HashToCrack matched.txt -passpol                     # With policy compliance
HashToCrack matched.txt -passpol -report             # Redacted for sharing
HashToCrack matched.txt -disabled -machines -passpol -report -o report.txt
```

**Sample Report:**
```
╔══════════════════════════════════════════════════════════════╗
║            HASHTOCRACK - Password Analytics Report           ║
╚══════════════════════════════════════════════════════════════╝

Filters Applied:
  • Disabled accounts: Excluded
  • Machine accounts:  Excluded

═══════════════════════════════════════════════════════════════
                      GENERAL STATISTICS                        
═══════════════════════════════════════════════════════════════

  Total Accounts Analyzed:  1500
  Passwords Cracked:        892 (59.47%)
  Passwords Not Cracked:    608 (40.53%)

  Crack Progress: [████████████████████████░░░░░░░░░░░░░░░░] 59.5%

═══════════════════════════════════════════════════════════════
                  PASSWORD LENGTH DISTRIBUTION                  
═══════════════════════════════════════════════════════════════

   8 chars: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  245 (27.5%)
   9 chars: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓            178 (20.0%)
  10 chars: ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓               156 (17.5%)

═══════════════════════════════════════════════════════════════
                    TOP 10 MOST USED PASSWORDS                  
═══════════════════════════════════════════════════════════════

  Rank  Password                        Count
  ────  ──────────────────────────────  ─────
  #1    Pas**********                   45      (with -report flag)
  #2    Sum*******                      32
  #3    Wel****                         28

═══════════════════════════════════════════════════════════════
              PASSWORD POLICY COMPLIANCE ANALYSIS               
═══════════════════════════════════════════════════════════════

  Policy: DOMAIN_PASSWORD_COMPLEX
  Requirements:
    • Minimum 8 characters
    • At least 3 of 4 categories:
      - Uppercase letters (A-Z)
      - Lowercase letters (a-z)
      - Digits (0-9)
      - Special characters (!@#$%^&*...)

  Results:
    Compliant passwords:     654 (73.32% of cracked)
    Non-compliant passwords: 238 (26.68% of cracked)

  Compliance: [█████████████████████████████░░░░░░░░░░░] 73.3%
```

## Command Reference

| Command | Description |
|---------|-------------|
| `HashToCrack help` | Display help message |
| `HashToCrack version` | Display version |
| `HashToCrack <file>` | Auto-detect mode based on file content |

### All Flags

| Flag | Modes | Description |
|------|-------|-------------|
| `-disabled` | All | Include disabled accounts |
| `-machines` | All | Include machine accounts (ending with `$`) |
| `-passpol` | Analytics | Show password policy compliance |
| `-report` | Analytics | Redact passwords in output |
| `-o`, `-outfile` | All | Write output to specified file |

## File Formats

### NTDS File Format

Expected format from tools like `secretsdump.py`:

```
domain\username:RID:LMHash:NTHash::: (status=Enabled)
domain\username:RID:LMHash:NTHash::: (status=Disabled)
```

### Hashcat Potfile Format

Standard hashcat potfile format:

```
hash:password
b4b9b02e6f09a9bd760f388b67351e2b:Summer2024!
```

## Building

### Build for Current Platform

```bash
go build -o HashToCrack ./cmd/HashToCrack
```

### Cross-Compile

```bash
# Using Makefile
make all        # All platforms
make windows    # Windows only
make linux      # Linux only
make darwin     # macOS only
make clean      # Clean build artifacts

# Manual cross-compilation
GOOS=linux GOARCH=amd64 go build -o HashToCrack_linux_amd64 ./cmd/HashToCrack
GOOS=darwin GOARCH=arm64 go build -o HashToCrack_darwin_arm64 ./cmd/HashToCrack
GOOS=windows GOARCH=amd64 go build -o HashToCrack_windows_amd64.exe ./cmd/HashToCrack
```

## Project Structure

```
HashToCrack/
├── cmd/
│   └── hashtocrack/
│       ├── main.go          # Entry point
│       └── version.go       # Version constant
├── internal/
│   ├── cli/
│   │   ├── cli.go           # Argument parsing
│   │   └── help.go          # Help messages
│   ├── ntds/
│   │   ├── types.go         # Data structures
│   │   ├── parser.go        # NTDS parsing
│   │   └── potfile.go       # Potfile loading
│   ├── modes/
│   │   ├── extract.go       # Extract mode
│   │   ├── match.go         # Match mode
│   │   └── analytics.go     # Analytics mode
│   └── utils/
│       └── utils.go         # Utilities
├── .github/
│   └── workflows/
│       ├── ci.yml           # CI pipeline
│       └── release.yml      # Release automation
├── go.mod
├── Makefile
├── LICENSE
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

⚠️ **This tool is intended for authorized security assessments only.**

Always ensure you have proper authorization before analyzing password data. The authors are not responsible for misuse of this tool.

## Acknowledgments

- Inspired by the need for better NTDS analysis tools
- Thanks to the security community for feedback and suggestions
