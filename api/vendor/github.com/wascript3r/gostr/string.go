package gostr

import (
	"unicode"
)

// UpperFirst capitalises the first character of a string
func UpperFirst(str string) string {
	if len(str) == 0 {
		return ""
	}
	tmp := []rune(str)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}
