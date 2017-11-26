package util

import "strings"

//LowerFirst LowerFirst
func LowerFirst(s string) string {
	if s == "" {
		return s
	}
	first := s[:1]
	return strings.ToLower(first) + s[1:]
}
