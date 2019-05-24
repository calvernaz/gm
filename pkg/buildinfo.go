package buildinfo


// GitInfo is either the empty string (the default)
// or is set to the git hash of the most recent commit
var GitInfo string

// Version is a string like "0.10" or "1.0", if applicable.
var Version string

// Summary returns the version and/or git version of this binary.
// If the linker flags were not provided, the return value is "unknown".
func Summary() string {
	if Version != "" && GitInfo != "" {
		return Version + ", " + GitInfo
	}
	if GitInfo != "" {
		return GitInfo
	}
	if Version != "" {
		return Version
	}
	return "unknown"
}
