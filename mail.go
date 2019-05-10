package sendmail

import (
	"errors"
	"io/ioutil"
)

// Mail containes an envelope and a content
type Mail struct {
	recipients       map[string]string
	carbonCopys      map[string]string
	blindCarbonCopys map[string]string
	attachments      map[string]string
	headers          map[string]string
	content          string
	isHTML           bool
}

// NewMail returns new mail
func NewMail() *Mail {
	m := &Mail{
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		"",
		false,
	}
	m.headers["MIME-Version"] = "1.0"
	m.headers["X-Mailer"] = "flagUpDown [sendmail]"
	return m
}

// SetFromEmail set sender
func (m *Mail) SetFromEmail(email, name string) {
	m.headers["Reply-To"] = email
	m.headers["From"] = contactEmailName(email, name)
}

// AddRecipient adds recipient
func (m *Mail) AddRecipient(email, name string) {
	m.recipients[email] = name
}

// AddCarbonCopy adds carbon copy recipient
func (m *Mail) AddCarbonCopy(email, name string) {
	m.carbonCopys[email] = name
}

// AddBlindCarbonCopy adds blind carbon copy recipient
func (m *Mail) AddBlindCarbonCopy(email, name string) {
	m.blindCarbonCopys[email] = name
}

// SetSubject sets subject
func (m *Mail) SetSubject(subject string) {
	m.headers["Subject"] = subject
}

// SetContent sets mail content
func (m *Mail) SetContent(content string, isHTML bool) {
	m.isHTML = isHTML
	m.content = content
}

// AddAttachment adds mail attachment by file path
func (m *Mail) AddAttachment(path, name string) {
	m.attachments[name] = path
}

func (m *Mail) getAllRcpt() ([]string, error) {
	length := len(m.recipients) + len(m.carbonCopys) + len(m.blindCarbonCopys)
	emailList := make([]string, length)
	err := errors.New("sendmail: recipient is empty")
	for email := range m.recipients {
		length--
		emailList[length] = email
		err = nil
	}
	for email := range m.carbonCopys {
		length--
		emailList[length] = email
	}
	for email := range m.blindCarbonCopys {
		length--
		emailList[length] = email
	}
	return emailList, err
}

func (m *Mail) setHeader(field, value string) {
	m.headers[field] = value
}

func (m *Mail) createPlainBody() string {
	body := ""
	body += "Content-Type: text/plain; charset=\"UTF-8\"" + CRLF
	body += "Content-Transfer-Encoding: base64" + CRLF + CRLF
	body += base64Encode([]byte(m.content)) + CRLF
	return body
}

func (m *Mail) createHTMLBody() string {
	const boundary = "_boundary_by_alternative_sendmail_"
	const startBoundary = "--" + boundary + CRLF
	const endBoundary = "--" + boundary + "--" + CRLF
	body := ""
	body += "Content-Type: multipart/alternative; boundary=\"" + boundary + "\"" + CRLF + CRLF
	body += startBoundary
	body += m.createPlainBody() + CRLF
	body += startBoundary
	body += "Content-Type: text/html; charset=\"UTF-8\"" + CRLF
	body += "Content-Transfer-Encoding: base64" + CRLF + CRLF
	body += base64Encode([]byte(m.content)) + CRLF + CRLF
	body += endBoundary
	return body
}

func (m *Mail) createBodyWithAttachment() string {
	const boundary = "_boundary_by_mixed_sendmail_"
	const startBoundary = "--" + boundary + CRLF
	const endBoundary = "--" + boundary + "--" + CRLF
	body := ""
	body += "Content-Type: multipart/mixed; boundary=\"" + boundary + "\"" + CRLF + CRLF
	body += startBoundary
	if m.isHTML {
		body += m.createHTMLBody()
	} else {
		body += m.createPlainBody()
	}
	for name, path := range m.attachments {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		body += CRLF
		body += startBoundary
		body += "Content-Type: application/octet-stream; name=\"" + name + "\"" + CRLF
		body += "Content-Transfer-Encoding: base64" + CRLF
		body += "Content-Disposition: attachment; filename=\"" + name + "\"" + CRLF + CRLF
		body += base64Encode(data) + CRLF + CRLF
	}
	body += endBoundary
	return body
}

func (m *Mail) createBody() string {
	if len(m.attachments) != 0 {
		return m.createBodyWithAttachment()
	} else if m.isHTML {
		return m.createHTMLBody()
	} else {
		return m.createPlainBody()
	}
}

func (m *Mail) toString() string {
	m.headers["To"] = mergeEmails(m.recipients)
	m.headers["Cc"] = mergeEmails(m.carbonCopys)
	message := ""
	for k, v := range m.headers {
		message += k + ": " + v + CRLF
	}

	message += m.createBody()

	return message
}
