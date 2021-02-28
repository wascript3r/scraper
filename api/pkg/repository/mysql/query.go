package mysql

import (
	"strconv"
	"strings"
)

const (
	SQLValuesParam = ":values"
)

type Row interface {
	Scan(dest ...interface{}) error
}

func GenerateParamString(start, n int) string {
	p := make([]string, n)
	for i := range p {
		p[i] = "$" + strconv.Itoa(start+i)
	}
	return strings.Join(p, ",")
}

func GenerateInsertValuesString(start, lines, values int) string {
	p := make([]string, lines)
	for i := range p {
		p[i] = "(" + GenerateParamString(start, values) + ")"
		start += values
	}
	return strings.Join(p, ",")
}

func ReplaceQueryValues(start, n int, query string) string {
	return strings.Replace(
		query,
		SQLValuesParam,
		GenerateParamString(start, n),
		1,
	)
}

func ReplaceQueryInsertValues(start, lines, values int, query string) string {
	return strings.Replace(
		query,
		SQLValuesParam,
		GenerateInsertValuesString(start, lines, values),
		1,
	)
}

func EncodeIDs(ids []int) []interface{} {
	a := make([]interface{}, len(ids))
	for i, v := range ids {
		a[i] = v
	}
	return a
}
