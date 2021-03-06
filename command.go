package sendmail

import (
	"crypto/tls"
	"errors"
	"io"
	"net/textproto"
	"strings"
)

func (c *Client) hello() error {
	err := c.ehlo()
	if err != nil {
		err = c.helo()
	}
	return err
}

func (c *Client) ehlo() error {
	_, msg, err := c.cmd(250, "EHLO %s", c.local)
	if err != nil {
		return err
	}
	extList := strings.Split(msg, "\n")
	if len(extList) > 1 {
		extList = extList[1:]
		for _, line := range extList {
			args := strings.SplitN(line, " ", 2)
			if len(args) > 1 {
				c.ext[args[0]] = args[1]
			} else {
				c.ext[args[0]] = ""
			}
		}
	}
	return err
}

func (c *Client) helo() error {
	c.ext = nil
	_, _, err := c.cmd(250, "HELO %s", c.local)
	return err
}

func (c *Client) plainAuth(identity, user, token string) error {
	resp := identity + "\x00" + user + "\x00" + token
	resp = base64Encode([]byte(resp))
	_, _, err := c.cmd(235, "AUTH PLAIN %s", resp)
	return err
}

func (c *Client) auth() error {
	var err error
	switch {
	case strings.ContainsAny(c.ext["AUTH"], "LOGIN"):
		err = c.loginAuth(c.authInfo.user, c.authInfo.password)
	case strings.ContainsAny(c.ext["AUTH"], "PLAIN"):
		err = c.plainAuth("", c.authInfo.user, c.authInfo.password)
	default:
		return errors.New("sendmail: auth failed")
	}
	return err
}

func (c *Client) loginAuth(user, password string) error {
	_, _, err := c.cmd(334, "AUTH LOGIN")
	if err != nil {
		goto RET
	}
	user = base64Encode([]byte(user))
	_, _, err = c.cmd(334, user)
	if err != nil {
		goto RET
	}
	password = base64Encode([]byte(password))
	_, _, err = c.cmd(235, password)
RET:
	return err
}

func (c *Client) mail() error {
	cmdStr := "MAIL FROM:<%s>"
	if c.ext != nil {
		if _, ok := c.ext["8BITMIME"]; ok {
			cmdStr += " BODY=8BITMIME"
		}
	}
	_, _, err := c.cmd(250, cmdStr, c.authInfo.user)
	return err
}

func (c *Client) rcpt(to string) error {
	if !validateLine(to) {
		return errors.New("smtp: recipient is not valid")
	}
	_, _, err := c.cmd(25, "RCPT TO:<%s>", to)
	return err
}

func (c *Client) data(data string) error {
	_, _, err := c.cmd(354, "DATA")
	if err != nil {
		goto RET
	}
	data += CRLF + "."
	_, _, err = c.cmd(0, data)
RET:
	return err
}

func (c *Client) startTLS(config *tls.Config) error {
	err := c.hello()
	if err != nil {
		goto RET
	}

	if _, _, err = c.cmd(220, "STARTTLS"); err != nil {
		goto RET
	}
	c.nativeConn = tls.Client(c.nativeConn, config)
	c.conn = textproto.NewConn(c.nativeConn)
	err = c.hello()
RET:
	return err
}

func (c *Client) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	var code int
	var msg string
CMD:
	id, err := c.conn.Cmd(format, args...)
	if err != nil {
		goto RET
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	code, msg, err = c.conn.ReadResponse(expectCode)
	if err == io.EOF {
		err = c.ReDial()
		if err != nil {
			goto RET
		}
		goto CMD
	} else if err != nil {
		goto RET
	}

RET:
	return code, msg, err
}
