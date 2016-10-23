package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pborman/uuid"
	"github.com/urfave/cli"
	"gopkg.in/gomail.v2"
)

func main() {
	app := cli.NewApp()
	app.Action = appMain
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "from, f",
			Usage: "email sender",
		},
		cli.StringSliceFlag{
			Name:  "recipient, r",
			Usage: "recipient(s) of email",
		},
		cli.StringFlag{
			Name:  "subject, s",
			Usage: "email subject",
			Value: "Your report is attached",
		},
		cli.StringSliceFlag{
			Name:  "attachment, a",
			Usage: "file(s) to attach to email",
		},
		cli.StringFlag{
			Name:  "mailqueue, q",
			Usage: "email mailqueue folder",
		},
	}
	app.HideVersion = true
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func appMain(c *cli.Context) error {
	from := c.String("from")
	to := c.StringSlice("recipient")
	subj := c.String("subject")
	attach := c.StringSlice("attachment")
	queue := c.String("mailqueue")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subj)
	m.SetBody("text/plain", "Hello,\n\nYour report is attached.\n\n")

	for _, file := range attach {
		m.Attach(file)
	}

	s := gomail.SendFunc(func(from string, to []string, msg io.WriterTo) error {
		p := filepath.Join(queue, uuid.New()+".eml")
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer func(x string) {
			if e := f.Close(); e != nil {
				log.Fatalf("closing file %q: %v", x, e)
			}
		}(p)

		if _, err = msg.WriteTo(f); err != nil {
			return err
		}
		return nil
	})

	return gomail.Send(s, m)
}
