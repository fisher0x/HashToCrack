package cli

import "fmt"

// PrintUsage displays a brief usage message
func PrintUsage() {
	fmt.Println(`HashToCrack - NTDS Hash Analyzer & Password Statistics Tool

Usage:
  HashToCrack <ntdsfile> [-disabled] [-machines] [-o <outfile>]
  HashToCrack <ntdsfile> <crackfile> [-disabled] [-machines] [-o <outfile>]
  HashToCrack <analyticsfile> [-disabled] [-machines] [-passpol] [-report] [-o <outfile>]
  HashToCrack help

Run 'HashToCrack help' for more information.`)
}

// PrintHelp displays the full help message
func PrintHelp(version string) {
	fmt.Printf(`HashToCrack v%s - NTDS Hash Analyzer & Password Statistics Tool
@fisher0x

DESCRIPTION:
  HashToCrack is a cross-platform tool for analyzing NTDS (Active Directory) 
  parsed files and hashcat potfiles. It helps security professionals match 
  cracked hashes with their owners and generate comprehensive statistics.

USAGE:
  HashToCrack <ntdsfile> [options]
  HashToCrack <ntdsfile> <crackfile> [options]
  HashToCrack <analyticsfile> [options]
  HashToCrack help

MODES:

  1. EXTRACT MODE - Extract hashes from NTDS file
     HashToCrack <ntdsfile> [-disabled] [-machines] [-o <outfile>]
     
     Extracts NT hashes from the NTDS file. By default, only enabled 
     user accounts are included (machine accounts excluded).
     
     Examples:
       HashToCrack NTDS.dit                    # Extract enabled user hashes
       HashToCrack NTDS.dit -disabled          # Include disabled accounts
       HashToCrack NTDS.dit -machines          # Include machine accounts
       HashToCrack NTDS.dit -o hashes.txt      # Save to file

  2. MATCH MODE - Match NTDS with cracked passwords
     HashToCrack <ntdsfile> <crackfile> [-disabled] [-machines] [-o <outfile>]
     
     Matches hashes from NTDS file with a hashcat potfile and displays
     usernames with their cracked passwords.
     
     Output format: username:hash:password:status
     
     Examples:
       HashToCrack NTDS.dit potfile.txt
       HashToCrack NTDS.dit potfile.txt -disabled -machines -o matched.txt

  3. ANALYTICS MODE - Generate password statistics
     HashToCrack <analyticsfile> [-disabled] [-machines] [-passpol] [-report] [-o <outfile>]
     
     Analyzes a matched file (output from match mode) and generates 
     comprehensive password statistics.
     
     Statistics include:
       - Total and cracked password counts
       - Password length distribution
       - Top 10 most common passwords
       - Password complexity compliance (DOMAIN_PASSWORD_COMPLEX)
     
     Examples:
       HashToCrack matched.txt -passpol
       HashToCrack matched.txt -disabled -machines -passpol
       HashToCrack matched.txt -passpol -report      # Redact passwords in output

OPTIONS:
  -disabled       Include disabled accounts in the analysis
  -machines       Include machine accounts (accounts ending with $)
  -passpol        Show password policy compliance statistics
  -report         Redact passwords in output (show first 3 chars only)
  -o, -outfile    Write output to specified file instead of stdout

NTDS FILE FORMAT:
  Expected format (secretsdump output):
  domain\username:RID:LMHash:NTHash::: (status=Enabled|Disabled)

CRACKFILE FORMAT:
  Standard hashcat potfile format:
  hash:password
`, version)
}
