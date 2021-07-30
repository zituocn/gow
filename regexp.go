/*
regexp.go

like mux  https://github.com/gorilla/mux
@see https://github.com/gorilla/mux/blob/master/regexp.go

*/

package gow

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// routeRegexp route regexp struct
type routeRegexp struct {
	path   string
	regexp *regexp.Regexp

	varsN []string
	varsR []*regexp.Regexp
}

// Match implements interface
//	search route
func (rr *routeRegexp) Match(path string, match *matchValue) bool {
	// set variables.
	rr.setMatch(path, match)

	//return true
	return rr.regexp.MatchString(path)
}

// setMatch return match.params
func (rr *routeRegexp) setMatch(path string, match *matchValue) {
	if rr.regexp != nil {
		matches := rr.regexp.FindStringSubmatchIndex(path)
		if len(matches) > 0 {
			vars := make(map[string]string)
			extractVars(path, matches, rr.varsN, vars)

			//vars to match.param
			params := new(Params)
			for k, v := range vars {
				*params = append(*params, Param{
					Key:   k,
					Value: v,
				})
			}
			match.params = params
		}
	}
}

// extractVars get vars
func extractVars(input string, matches []int, names []string, output map[string]string) {
	for i, name := range names {
		output[name] = input[matches[2*i+2]:matches[2*i+3]]
	}
}

// getPattern return pattern
func getPattern(path string) (pattern string) {
	pattern = "\\w+"
	if strings.Contains(path, "{static_file_path}") || strings.Contains(path, "{match_all}") {
		pattern = ".*"
	}
	return pattern
}

// addRouteRegexp add regexp route to math
func addRouteRegexp(path string, rc *routeConfig) (*routeRegexp, error) {
	path = cleanPath(path)
	idxs, errBraces := braceIndices(path)
	if errBraces != nil {
		return nil, errBraces
	}
	fullPath := path
	defaultPattern := getPattern(path)
	varsN := make([]string, len(idxs)/2)
	varsR := make([]*regexp.Regexp, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	// if ignoreCase is true
	if rc.ignoreCase {
		pattern.WriteString("(?i)")
	}
	pattern.WriteByte('^')
	reverse := bytes.NewBufferString("")
	var end int
	var err error
	for i := 0; i < len(idxs); i += 2 {
		// Set all values we are interested in.
		raw := path[end:idxs[i]]
		end = idxs[i+1]
		parts := strings.SplitN(path[idxs[i]+1:end-1], ":", 2)
		name := parts[0]
		patt := defaultPattern
		if len(parts) == 2 {
			patt = parts[1]
		}
		// Name or pattern can't be empty.
		if name == "" || patt == "" {
			return nil, fmt.Errorf("mux: missing name or pattern in %q", path[idxs[i]:end])
		}
		// Build the regexp pattern.
		fmt.Fprintf(pattern, "%s(%s)", regexp.QuoteMeta(raw), patt)

		// Build the reverse template.
		fmt.Fprintf(reverse, "%s%%s", raw)
		// Append variable name and compiled pattern.
		varsN[i/2] = name
		varsR[i/2], err = regexp.Compile(fmt.Sprintf("^%s$", patt))
		if err != nil {
			return nil, err
		}
	}
	// Add the remaining.
	raw := path[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	pattern.WriteByte('$')
	reg, errCompile := regexp.Compile(pattern.String())
	if errCompile != nil {
		return nil, errCompile
	}

	if reg.NumSubexp() != len(idxs)/2 {
		panic(fmt.Sprintf("route %s contains capture groups in its regexp. ", fullPath) +
			"Only non-capturing groups are accepted: e.g. (?:pattern) instead of (pattern)")
	}

	return &routeRegexp{
		path:   fullPath,
		regexp: reg,
		varsR:  varsR,
		varsN:  varsN,
	}, nil
}

func braceIndices(s string) ([]int, error) {
	var level, idx int
	var idxs []int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("mux: unbalanced braces in %q", s)
	}
	return idxs, nil
}
