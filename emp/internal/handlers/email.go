// internal/handlers/email.go

package handlers

import (
	"fmt"
	"io"
	_ "log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/go-gomail/gomail"
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/template/html/v2"
)

const (
	imapServer = "imap.gmail.com:993"
	smtpServer = "smtp.gmail.com"
	smtpPort   = 587
	email      = "dchildir@gmail.com"
	password   = "muzl xgpe cpik kdek"
)

type Email struct {
	From    string
	Subject string
	Date    time.Time
	Body    string
}

func HandleIndex(c *fiber.Ctx) error {
	return c.Render("index2", fiber.Map{})
}

func HandleGetEmails(c *fiber.Ctx) error {
	emails, err := fetchEmails()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("email", fiber.Map{"Emails": emails})
}

func HandleSendEmail(c *fiber.Ctx) error {
	to := c.FormValue("to")
	subject := c.FormValue("subject")
	body := c.FormValue("body")

	err := sendEmail(to, subject, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("Email sent successfully")
}

func fetchEmails() ([]Email, error) {
	c, err := client.DialTLS(imapServer, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}
	defer c.Logout()

	if err := c.Login(email, password); err != nil {
		return nil, fmt.Errorf("error logging in: %v", err)
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("error selecting INBOX: %v", err)
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 5 {
		from = to - 4
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var emails []Email

	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			continue
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			continue
		}

		header := mr.Header
		date, _ := header.Date()
		from, _ := header.AddressList("From")
		subject, _ := header.Subject()

		var body strings.Builder
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				break
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				b, _ := io.ReadAll(p.Body)
				body.Write(b)
			case *mail.AttachmentHeader:
				// Handle attachments if needed
			}
		}

		emails = append(emails, Email{
			From:    from[0].String(),
			Subject: subject,
			Date:    date,
			Body:    body.String(),
		})
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}

	return emails, nil
}

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpServer, smtpPort, email, password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
