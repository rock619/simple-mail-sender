package main

import (
	"flag"
	"fmt"
	"log"

	sms "github.com/rock619/simple-mail-sender"
)

func main() {
	if err := newExec().Do(); err != nil {
		log.Fatal(err)
	}
}

type Exec struct {
	from, to, subject, text string
}

func newExec() *Exec {
	e := &Exec{}

	flag.StringVar(&e.from, "from", "", "from address")
	flag.StringVar(&e.to, "to", "", "to address")
	flag.StringVar(&e.subject, "subject", "", "subject")
	flag.StringVar(&e.text, "text", "", "body text")
	flag.Parse()

	return e
}

func (e *Exec) Do() error {
	host, err := sms.ResolveMXHost(e.to)
	if err != nil {
		return fmt.Errorf("sms.ResolveMXHost: %w", err)
	}

	sender, err := sms.NewSMTPSender(host)
	if err != nil {
		return fmt.Errorf("sms.NewSMTPSender(%s): %w", host, err)
	}

	err = sms.SendMail(sender, e.from, e.to, e.subject, e.text)
	if err != nil {
		return fmt.Errorf("sms.SendMail: %w", err)
	}
	return nil
}
