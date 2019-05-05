package sendmail

// Mail containes an envelope and a content
type Mail struct {
	headers map[string]string
	content string
}

// NewMail returns new mail
func NewMail() *Mail {
	return &Mail{make(map[string]string), ""}
}

// SetRecipient sets recipient
func (m *Mail) SetRecipient(recipient string) {
	m.headers["to"] = recipient
}

// SetSubject sets subject
func (m *Mail) SetSubject(subject string) {
	m.headers["subject"] = subject
}

// SetContent sets mail content
func (m *Mail) SetContent(content string) {
	m.content = content
}

func (m *Mail) setHeader(field, value string) {
	m.headers[field] = value
}

func (m *Mail) toString() string {
	message := ""
	for k, v := range m.headers {
		message += k + ":" + v + "\r\n"
	}
	message += "\r\n" + m.content
	return message
}
