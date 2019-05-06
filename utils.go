package sendmail

import (
	"fmt"
	"strings"
)

// CRLF represents Carriage-Return Line-Feed
const CRLF = "\r\n"

func validateLine(line string) bool {
	flag := true
	if line == "" || strings.ContainsAny(line, "\n\r") {
		flag = false
	}
	return flag
}

func addr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
