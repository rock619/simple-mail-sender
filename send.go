package sms

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jhillyerd/enmime"
)

func SendMail(s enmime.Sender, from, to, subject, text string) error {
	err := enmime.Builder().
		From("", from).
		Subject(subject).
		Text([]byte(text)).
		To("", to).
		Send(s)
	if err != nil {
		return fmt.Errorf("enmime.MailBuilder.Send(): %w", err)
	}
	return nil
}

type SMTPSender struct {
	c *smtp.Client
}

func NewSMTPSender(host string) (enmime.Sender, error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", host, 25))
	if err != nil {
		return nil, fmt.Errorf("smtp.Dial: %w", err)
	}

	return &SMTPSender{
		c,
	}, nil
}

func (s *SMTPSender) Send(reversePath string, recipients []string, msg []byte) error {
	if err := s.c.Mail(reversePath); err != nil {
		return fmt.Errorf("*smtp.Client.Mail(%s): %w", reversePath, err)
	}
	for _, r := range recipients {
		if err := s.c.Rcpt(r); err != nil {
			return fmt.Errorf("*smtp.Client.Rcpt(%s): %w", r, err)
		}
	}

	wc, err := s.c.Data()
	if err != nil {
		return fmt.Errorf("*smtp.Client.Data(): %w", err)
	}
	if _, err := wc.Write(msg); err != nil {
		return fmt.Errorf("io.WriteCloser.Write(): %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("io.WriteCloser.Close(): %w", err)
	}

	if err := s.c.Quit(); err != nil {
		return fmt.Errorf("*smtp.Client.Quit(): %w", err)
	}
	return nil
}

type LogSender struct{}

func (s *LogSender) Send(reversePath string, recipients []string, msg []byte) error {
	log.Printf("reversePath=%s recipients=%s msg=%s", reversePath, recipients, msg)
	return nil
}
