package tmuxcli

import (
	"context"
	"strconv"
	"strings"
)

// -V            Report the tmux version.

// Version is a tmux version such as "3.7c".
type Version struct {
	// Raw is the version as tmux reports it, without the leading "tmux ".
	Raw    string
	Major  int
	Minor  int
	Suffix string // e.g. "c" in "3.7c"
}

// Version runs tmux -V.
// A version not of the form <major>.<minor><suffix>, optionally prefixed
// with "next-", is returned with only Raw set.
func (s *ShellOut) Version(ctx context.Context) (Version, error) {
	out, err := s.run(ctx, "-V")
	if err != nil {
		return Version{}, err
	}
	return parseVersion(out), nil
}

// Less reports whether v is older than other.
// Suffixes compare lexically, so "3.7" < "3.7a" < "3.7c".
func (v Version) Less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Suffix < other.Suffix
}

func parseVersion(out string) Version {
	raw := strings.TrimPrefix(strings.TrimSpace(out), "tmux ")
	v := Version{Raw: raw}
	majorStr, rest, ok := strings.Cut(strings.TrimPrefix(raw, "next-"), ".")
	if !ok {
		return v
	}
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return v
	}
	digits := len(rest) - len(strings.TrimLeft(rest, "0123456789"))
	minor, err := strconv.Atoi(rest[:digits])
	if err != nil {
		return v
	}
	v.Major, v.Minor, v.Suffix = major, minor, rest[digits:]
	return v
}
