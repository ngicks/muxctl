package zellijcli

import (
	"context"
	"strconv"
	"strings"
)

//   -V, --version
//           Print version

// Version is the version zellij reports, e.g. "zellij 0.45.1".
type Version struct {
	Raw   string // the reported line, e.g. "zellij 0.45.1"
	Major int
	Minor int
	Patch int
}

// Version runs zellij --version.
// Major, Minor and Patch stay 0 when the reported version is not "<int>.<int>.<int>",
// ignoring a pre-release or build suffix.
func (s *ShellOut) Version(ctx context.Context) (Version, error) {
	out, err := s.run(ctx, "--version")
	if err != nil {
		return Version{}, err
	}
	return parseVersion(out), nil
}

func parseVersion(out string) Version {
	raw := strings.TrimSpace(out)
	num := strings.TrimPrefix(raw, "zellij ")
	if i := strings.IndexAny(num, "-+"); i >= 0 {
		num = num[:i]
	}
	parts := strings.Split(num, ".")
	if len(parts) != 3 {
		return Version{Raw: raw}
	}
	var n [3]int
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return Version{Raw: raw}
		}
		n[i] = v
	}
	return Version{Raw: raw, Major: n[0], Minor: n[1], Patch: n[2]}
}
