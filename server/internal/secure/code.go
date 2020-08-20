package secure

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"math/rand"
	"time"
)

// An interface for types that send verification codes to a destination.
type CodeSender interface {
	SendCode(destination string) (code string, err error)
}

type codeGenerator struct {
	length  int
	charset string
}

func (cg codeGenerator) generateCode() string {
	b := make([]byte, cg.length)
	for i := range b {
		b[i] = cg.charset[rand.Intn(len(cg.charset))]
	}
	return string(b)
}

type emailSender struct {
	client        *sendgrid.Client
	codeGenerator codeGenerator
	senderEmail   *mail.Email
	subject       string
	content       string
	htmlContent   string
}

func EmailSender(apiKey string) CodeSender {
	rand.Seed(time.Now().UTC().UnixNano())
	return &emailSender{
		client: sendgrid.NewSendClient(apiKey),
		codeGenerator: codeGenerator{
			length:  6,
			charset: "0123456789",
		},
		senderEmail: mail.NewEmail("Kript", "verify@kript.us"),
		subject:     "Your Verification Code for Kript",
		content:     "Your verification code is %s.",
		htmlContent: "Your verification code is <b>%s</b>.",
	}
}

func (s *emailSender) SendCode(destination string) (code string, err error) {
	code = s.codeGenerator.generateCode()
	message := mail.NewSingleEmail(
		s.senderEmail,
		s.subject,
		mail.NewEmail("", destination),
		fmt.Sprintf(s.content, code),
		fmt.Sprintf(s.htmlContent, code))

	_, err = s.client.Send(message)
	if err != nil {
		code = ""
	}
	return
}
