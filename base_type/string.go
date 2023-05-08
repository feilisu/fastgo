package base_type

import "strings"

func FirstToLow(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
