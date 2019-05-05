package sendmail

import (
	"net/textproto"
)

type auth struct {
	identity string
	user     string
	password string
}

// Client represents a client connection to an SMTP server
type Client struct {
	server   string
	local    string
	conn     *textproto.Conn
	ext      map[string]string
	authInfo *auth
}

// Dial connects an SMTP server
func Dial(addr string) (*Client, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if _, _, err = conn.ReadResponse(220); err != nil {
		conn.Close()
		return nil, err
	}

	c := &Client{
		server:   addr,
		local:    "localhost",
		conn:     conn,
		ext:      make(map[string]string),
		authInfo: &auth{},
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
	c.authInfo = &auth{"", user, token}
}

// Send sends mail
func (c *Client) Send(mail *Mail) {
	c.hello()
	c.auth()
	c.mail()
	c.rcpt(mail.headers["to"])
	c.data(mail.toString())
}
