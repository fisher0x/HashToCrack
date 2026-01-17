package ntds

// Entry represents a parsed NTDS entry
type Entry struct {
	Username   string
	RID        string
	LMHash     string
	NTHash     string
	IsDisabled bool
	IsMachine  bool
	RawLine    string
}

// CrackedEntry represents a matched entry with password
type CrackedEntry struct {
	Entry
	Password string
	Cracked  bool
}

// AnalyticsResult holds statistics about cracked passwords
type AnalyticsResult struct {
	TotalAccounts      int
	CrackedAccounts    int
	CrackPercentage    float64
	LengthDistribution map[int]int
	TopPasswords       []PasswordCount
	ComplexCount       int
	ComplexPercentage  float64
}

// PasswordCount for top passwords ranking
type PasswordCount struct {
	Password string
	Count    int
}
