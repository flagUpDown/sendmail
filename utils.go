package sendmail

import (
	"encoding/base64"
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

func utf8B(str string) string {
	return fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(str)))
}

func contactEmailName(email, name string) string {
	if name != "" {
		name = utf8B(name) + " "
	}
	return name + "<" + email + ">"
}

func mergeEmails(emails map[string]string) string {
	result := ""
	for email, name := range emails {
		result += contactEmailName(email, name) + ", "
	}
	if result != "" {
		result = result[0 : len(result)-2]
	}
	return result
}
