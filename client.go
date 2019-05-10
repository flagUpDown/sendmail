package sendmail

import (
	"crypto/tls"
	"net"
	"net/textproto"
)

type auth struct {
	identity string
	user     string
	password string
}

// Client represents a client connection to an SMTP server
type Client struct {
	serverHost string
	serverPort int
	local      string
	isTLS      bool
	conn       *textproto.Conn
	nativeConn net.Conn
	ext        map[string]string
	authInfo   *auth
}

// Dial connects an SMTP server
func Dial(host string, port int, isTLS bool) (*Client, error) {
	nativeConn, err := net.Dial("tcp", addr(host, port))
	if err != nil {
		return nil, err
	}

	if isTLS {
		config := &tls.Config{ServerName: host}
		nativeConn = tls.Client(nativeConn, config)
	}

	conn := textproto.NewConn(nativeConn)

	if _, _, err = conn.ReadResponse(220); err != nil {
		conn.Close()
		return nil, err
	}

	c := &Client{
		serverHost: host,
		serverPort: port,
		local:      "localhost",
		isTLS:      isTLS,
		nativeConn: nativeConn,
		conn:       conn,
		ext:        make(map[string]string),
		authInfo:   &auth{},
	}

	return c, nil
}

// ReDial reconnects an SMTP server
func (c *Client) ReDial() error {
	var err error
	c.nativeConn, err = net.Dial("tcp", addr(c.serverHost, c.serverPort))
	if err != nil {
		return err
	}

	if c.isTLS {
		config := &tls.Config{ServerName: c.serverHost}
		c.nativeConn = tls.Client(c.nativeConn, config)
	}

	c.conn = textproto.NewConn(c.nativeConn)

	if _, _, err = c.conn.ReadResponse(220); err != nil {
		c.conn.Close()
		return err
	}
	return err
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Quit sends the QUIT command and closes the connection to the server.
func (c *Client) Quit() error {
	_, _, err := c.cmd(221, "QUIT")
	if err != nil {
		return err
	}
	return c.Close()
}

// SetAuth sets authentication infomation
func (c *Client) SetAuth(user, token string) {
	c.authInfo.user = user
	c.authInfo.password = token
}

// Send sends mail
func (c *Client) Send(mail *Mail) error {
	var emailList []string
	err := c.hello()
	if err != nil {
		goto RET
	}
	if _, ok := c.ext["STARTTLS"]; ok && !c.isTLS {
		config := &tls.Config{ServerName: c.serverHost}
		err = c.startTLS(config)
		if err != nil {
			goto RET
		}
	}
	err = c.auth()
	if err != nil {
		goto RET
	}
	err = c.mail()
	if err != nil {
		goto RET
	}
	emailList, err = mail.getAllRcpt()
	if err != nil {
		goto RET
	}
	for _, email := range emailList {
		err = c.rcpt(email)
		if err != nil {
			goto RET
		}
	}
	err = c.data(mail.toString())
RET:
	return err
}
