package sendmail

import "strings"

func validateLine(line string) bool {
	flag := true
	if line == "" || strings.ContainsAny(line, "\n\r") {
		flag = false
	}
	return flag
}
