package engine

import "regexp"

// versionHint matches an engine hint appended to a version string, e.g.
// "1.6.0 (tofu)".
var versionHint = regexp.MustCompile(`^\s*(.*?)\s*\(\s*([^()]+?)\s*\)\s*$`)

// ParseVersion splits an optional engine hint off a version string, returning
// the bare version and the hinted engine. The engine is nil when the string
// carries no hint, in which case the caller decides which engine to use.
func ParseVersion(s string) (string, *Engine, error) {
	m := versionHint.FindStringSubmatch(s)
	if m == nil {
		return s, nil, nil
	}
	e, err := Lookup(m[2])
	if err != nil {
		return "", nil, err
	}
	return m[1], e, nil
}
