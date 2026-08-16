// Package globx matches path globs with ** support (entry tags, §3.2/§7.1).
package globx

import (
	"regexp"
	"strings"
	"sync"
)

var (
	mu    sync.Mutex
	cache = map[string]*regexp.Regexp{}
)

// compile turns a glob like "agentdeskd/**" or "*.go" into a regexp.
func compile(glob string) *regexp.Regexp {
	mu.Lock()
	defer mu.Unlock()
	if re, ok := cache[glob]; ok {
		return re
	}
	var b strings.Builder
	b.WriteString("^")
	i := 0
	for i < len(glob) {
		switch {
		case strings.HasPrefix(glob[i:], "**/"):
			b.WriteString(`(?:[^/]+/)*`)
			i += 3
		case strings.HasPrefix(glob[i:], "**"):
			b.WriteString(`.*`)
			i += 2
		case glob[i] == '*':
			b.WriteString(`[^/]*`)
			i++
		case glob[i] == '?':
			b.WriteString(`[^/]`)
			i++
		default:
			b.WriteString(regexp.QuoteMeta(string(glob[i])))
			i++
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		re = regexp.MustCompile("^$")
	}
	cache[glob] = re
	return re
}

// Match reports whether path matches glob. A glob with no slash and no meta
// also matches any path segment (tag "supervisor" matches "supervisor/x.go").
func Match(glob, path string) bool {
	glob = strings.TrimPrefix(glob, "./")
	path = strings.TrimPrefix(path, "./")
	if compile(glob).MatchString(path) {
		return true
	}
	// "dir/**" should also match "dir" itself.
	if s, ok := strings.CutSuffix(glob, "/**"); ok && s == path {
		return true
	}
	// Bare token: match as prefix directory or any segment.
	if !strings.ContainsAny(glob, "/*?") {
		if strings.HasPrefix(path, glob+"/") {
			return true
		}
		for _, seg := range strings.Split(path, "/") {
			if seg == glob {
				return true
			}
		}
	}
	return false
}

// AnyMatch reports whether any glob matches any path.
func AnyMatch(globs, paths []string) bool {
	for _, g := range globs {
		for _, p := range paths {
			if Match(g, p) {
				return true
			}
		}
	}
	return false
}

// Valid reports whether the glob compiles.
func Valid(glob string) bool {
	return compile(glob) != nil
}
