package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Client struct {
	host     string
	port     string
	email    string
	password string
}

func New(host, port, email, password string) *Client {
	return &Client{host: host, port: port, email: email, password: password}
}

func (c *Client) Send(to, subject, body string) error {
	if c.email == "" || c.password == "" {
		return fmt.Errorf("smtp credentials are not configured")
	}
	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	msg := strings.Join([]string{
		"From: " + c.email,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
	auth := smtp.PlainAuth("", c.email, c.password, c.host)
	return smtp.SendMail(addr, auth, c.email, []string{to}, []byte(msg))
}
