package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"time"

	"project_layout/model/config"

	mail "github.com/xhit/go-simple-mail/v2"
	"go.uber.org/zap"
)

func NewEmail(
	host config.Email,
	logger *zap.SugaredLogger,
) IEmail {
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = host.Host
	server.Port = host.Port
	server.Username = host.Username
	server.Password = string(host.Password)
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = true
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := server.Connect()
	if err != nil {
		panic(err)
	}

	return &emailImpl{
		logger:     logger,
		from:       host.Username,
		smtpClient: smtpClient,
	}
}

const htmlTemp = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>Dear {{.Dear}}</p>
		<p>&emsp;{{.Content}}</p>
	</body>
</html>`

type Message struct {
	To      string `json:"to"`
	Title   string `json:"title"`
	Dear    string `json:"dear"`
	Content string `json:"content"`
}

func (msg Message) Check() error {
	if msg.To == "" {
		return fmt.Errorf("to is empty")
	}
	if msg.Title == "" {
		return fmt.Errorf("title is empty")
	}
	if msg.Content == "" {
		return fmt.Errorf("content is empty")
	}

	return nil
}

type IEmail interface {
	Send(ctx context.Context, msg Message) error
}

type emailImpl struct {
	logger     *zap.SugaredLogger
	from       string
	smtpClient *mail.SMTPClient
}

func (impl *emailImpl) Send(ctx context.Context, msg Message) error {
	if err := msg.Check(); err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(impl.from).
		AddTo(msg.To).
		SetSubject(msg.Title)

	buf := new(bytes.Buffer)
	temp, err := template.New("email-body").Parse(htmlTemp)
	if err != nil {
		return err
	}
	if err := temp.Execute(buf, msg); err != nil {
		return err
	}
	email.SetBody(mail.TextHTML, buf.String())

	if email.Error != nil {
		return email.Error
	}

	if err := email.Send(impl.smtpClient); err != nil {
		return err
	}

	return nil
}
