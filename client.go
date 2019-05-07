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
	conn       *textproto.Conn
	nativeConn net.Conn
	ext        map[string]string
	authInfo   *auth
}

// Dial connects an SMTP server
func Dial(host string, port int) (*Client, error) {
	nativeConn, err := net.Dial("tcp", addr(host, port))
	if err != nil {
		return nil, err
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
		nativeConn: nativeConn,
		conn:       conn,
		ext:        make(map[string]string),
		authInfo:   &auth{},
	}

	return c, nil
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
func (c *Client) Send(mail *Mail) {
	c.hello()
	if _, ok := c.ext["STARTTLS"]; ok {
		config := &tls.Config{ServerName: c.serverHost}
		c.startTLS(config)
	}
	c.auth()
	c.mail()
	for email := range mail.recipients {
		c.rcpt(email)
	}
	for email := range mail.carbonCopys {
		c.rcpt(email)
	}
	for email := range mail.blindCarbonCopys {
		c.rcpt(email)
	}
	c.rcpt(mail.headers["to"])
	c.data(mail.toString())
}
