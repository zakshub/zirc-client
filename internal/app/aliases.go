package app

import (
	"fmt"
	"strings"
	"unicode"
)

// expandAlias supports positional $1..$9, the remaining arguments $1-,
// and all arguments $*. It deliberately does not evaluate expressions,
// access the filesystem, launch processes, or expand environment variables.
func expandAlias(template string, args []string) (string, error) {
	if !strings.HasPrefix(template, "/") {
		return "", fmt.Errorf("alias must expand to a command")
	}
	var b strings.Builder
	for i := 0; i < len(template); {
		if template[i] != '$' {
			b.WriteByte(template[i])
			i++
			continue
		}
		if i+1 < len(template) && template[i+1] == '$' {
			b.WriteByte('$')
			i += 2
			continue
		}
		if i+1 < len(template) && template[i+1] == '*' {
			b.WriteString(strings.Join(args, " "))
			i += 2
			continue
		}
		if i+1 < len(template) && template[i+1] >= '1' && template[i+1] <= '9' {
			n := int(template[i+1] - '0')
			i += 2
			if i < len(template) && template[i] == '-' {
				if n <= len(args) { b.WriteString(strings.Join(args[n-1:], " ")) }
				i++
			} else if n <= len(args) {
				b.WriteString(args[n-1])
			}
			continue
		}
		b.WriteByte(template[i])
		i++
	}
	return strings.TrimSpace(b.String()), nil
}

func validAliasName(name string) bool {
	if name == "" || !unicode.IsLetter(rune(name[0])) { return false }
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' { return false }
	}
	return true
}
