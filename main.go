package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pborman/uuid"
	"github.com/scorredoira/email"
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

	m := email.NewMessage(subject, "Hi;\n\nYour report is attached.\n\nThank you")
	m.From = sender
	m.To = []string{recipient}

	if err := m.Attach(reportFile); err != nil {
		log.Fatal(err)
	}

	data := m.Bytes()
	name := uuid.New()

	target := filepath.Join(mailqueue, name+".eml")
	if err := ioutil.WriteFile(target, data, 0644); err != nil {
		log.Fatal(err)
	}
}
