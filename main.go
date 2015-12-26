package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pborman/uuid"

	"gopkg.in/gomail.v2"
)

var (
	reportFile string
	recipient  string
	sender     string
	subject    string = "Your report is attached"
	mailqueue  string
)

func init() {
	flag.StringVar(&reportFile, "file", "", "file name of report")
	flag.StringVar(&recipient, "to", "", "recipient of email report")
	flag.StringVar(&sender, "from", "", "sender of email report")
	flag.StringVar(&subject, "subject", subject, "subject of email report")
	flag.StringVar(&mailqueue, "q", "", "mail queue folder")
}

func main() {
	flag.Parse()
	if reportFile == "" {
		log.Fatal("report file required")
	}
	if recipient == "" {
		log.Fatal("recipient required")
	}
	if sender == "" {
		log.Fatal("sender required")
	}
	if mailqueue == "" {
		log.Fatal("mailqueue required")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", "Hello,\n\nYour report is attached.\n\n")
	m.Attach(reportFile)

	s := gomail.SendFunc(func(from string, to []string, msg io.WriterTo) error {
		p := filepath.Join(mailqueue, uuid.New()+".eml")
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = msg.WriteTo(f); err != nil {
			return err
		}
		return nil
	})

	if err := gomail.Send(s, m); err != nil {
		log.Fatal(err)
	}
}
