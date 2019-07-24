package version

import "fmt"

var (
	version = "untagged"
	commit  string
	branch  string
)

// String returns a human-readable version string.
func String() string {
	hasVersion := version != ""
	hasBuildInfo := commit != ""

	switch {
	case hasVersion && hasBuildInfo:
		return fmt.Sprintf("version %s (commit %s, branch %s)", version, commit, branch)
	case !hasVersion && hasBuildInfo:
		return fmt.Sprintf("(commit %s, branch %s)", commit, branch)
	case hasVersion && !hasBuildInfo:
		return fmt.Sprintf("version %s (no build information)", version)
	default:
		return "(no version or build info)"
	}
}

// Version returns the version string.
func Version() string { return version }

// CommitHash returns the commit hash at which the binary was built.
func CommitHash() string { return commit }

// Branch returns the branch at which the binary was built.
func Branch() string { return branch }
