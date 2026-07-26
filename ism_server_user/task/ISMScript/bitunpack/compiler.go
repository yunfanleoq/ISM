package bitunpack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reBitGetAssign = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*BitGet\s*\(\s*"([^"]+)"\s*,\s*(\d+)\s*\)\s*$`)
	reSetDeviceVar = regexp.MustCompile(`^\s*SetDeviceData\s*\(\s*"([^"]+)"\s*,\s*([A-Za-z_][A-Za-z0-9_]*)\s*\)\s*$`)
	reAssignOnly   = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*$`)
)

type bitBinding struct {
	SourceKey string
	Bit       uint8
}

func stripLineComment(line string) string {
	if idx := strings.Index(line, "//"); idx >= 0 {
		return strings.TrimSpace(line[:idx])
	}
	return line
}

func isSkippable(line string) bool {
	return line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/*")
}

// normalizeLines joins "VAR =" with the following non-empty statement line.
func normalizeLines(content string) []string {
	raw := strings.Split(content, "\n")
	out := make([]string, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		line := stripLineComment(strings.TrimSpace(raw[i]))
		if isSkippable(line) {
			continue
		}
		if reAssignOnly.MatchString(line) {
			merged := false
			for j := i + 1; j < len(raw); j++ {
				next := stripLineComment(strings.TrimSpace(raw[j]))
				if isSkippable(next) {
					continue
				}
				line = strings.TrimSpace(line + " " + next)
				i = j
				merged = true
				break
			}
			if !merged {
				continue
			}
		}
		out = append(out, line)
	}
	return out
}

// Compile tries to compile a pure BitGet+SetDeviceData script into rules.
// Unknown SetDeviceData target vars (script typos) are skipped; unrecognized
// statement forms cause ok=false so the script falls back to Anko.
func Compile(scriptUUID, scriptName, content string) (rules []Rule, ok bool) {
	rules, _, ok = CompileWithError(scriptUUID, scriptName, content)
	return rules, ok
}

// CompileWithError is like Compile but returns the first rejected line for diagnostics.
func CompileWithError(scriptUUID, scriptName, content string) (rules []Rule, rejectLine string, ok bool) {
	vars := make(map[string]bitBinding)
	for _, line := range normalizeLines(content) {
		if m := reBitGetAssign.FindStringSubmatch(line); m != nil {
			src := strings.TrimSpace(m[2])
			bit64, err := strconv.ParseUint(m[3], 10, 8)
			// Skip malformed BitGet (no device->point, or bad bit): same as runtime BitGet failure.
			if err != nil || bit64 == 0 || !strings.Contains(src, "->") {
				continue
			}
			vars[m[1]] = bitBinding{SourceKey: src, Bit: uint8(bit64)}
			continue
		}
		if m := reSetDeviceVar.FindStringSubmatch(line); m != nil {
			target := strings.TrimSpace(m[1])
			parts := strings.SplitN(target, "->", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return nil, line, false
			}
			bind, exists := vars[m[2]]
			if !exists {
				// Likely a typo in the script (e.g. DCQY22 vs DCQY2); skip this rule.
				continue
			}
			srcParts := strings.SplitN(bind.SourceKey, "->", 2)
			if len(srcParts) != 2 {
				return nil, line, false
			}
			rules = append(rules, Rule{
				SourceDevice: srcParts[0],
				SourcePoint:  srcParts[1],
				Bit:          bind.Bit,
				TargetDevice: parts[0],
				TargetPoint:  parts[1],
				ScriptUUID:   scriptUUID,
				ScriptName:   scriptName,
			})
			continue
		}
		return nil, line, false
	}
	if len(rules) == 0 {
		return nil, "", false
	}
	return rules, "", true
}

// SourceKey returns "device->point" for a rule.
func (r Rule) SourceKey() string {
	return fmt.Sprintf("%s->%s", r.SourceDevice, r.SourcePoint)
}

// TargetKey returns "device->point" for a rule target.
func (r Rule) TargetKey() string {
	return fmt.Sprintf("%s->%s", r.TargetDevice, r.TargetPoint)
}
