package query_corrector

import "strings"

func Correct(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)
	return query
}
