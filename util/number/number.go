package number

import "regexp"

var isString = regexp.MustCompile(`\D`)

func ClearNumber(v string) string {
	return isString.ReplaceAllString(v, "")
}
