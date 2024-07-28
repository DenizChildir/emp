// cmd/main.go

package main

import (
	"emp/internal/database"
	"emp/internal/handlers"
	"emp/internal/utils"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/gofiber/template/html/v2"
	"io"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/go-gomail/gomail"
)

func main() {
	// Set up template engine
	utils.InitTwilio("YOUR_TWILIO_ACCOUNT_SID", "YOUR_TWILIO_AUTH_TOKEN")
	engine := html.New("./templates", ".html")
	engine.AddFunc("iterate", func(start, end int) []int {
		var result []int
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
		return result
	})
	engine.AddFunc("sum", func(data interface{}, field string) float64 {
		var total float64
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(data)
			for i := 0; i < s.Len(); i++ {
				total += s.Index(i).FieldByName(field).Float()
			}
		}
		return total
	})

	// Set up Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files
	app.Static("/", "./static")

	// Set up database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Set up routes
	setupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App) {

	app.Get("/home", handlers.Home)

	// Employee routes
	app.Get("/employees", handlers.GetEmployees)
	app.Get("/employees/new", handlers.NewEmployee)
	app.Post("/employees", handlers.CreateEmployee)
	app.Get("/employees/:id", handlers.GetEmployee)
	app.Get("/employees/:id/edit", handlers.EditEmployee)
	app.Post("/employees/:id", handlers.UpdateEmployee)
	app.Delete("/employees/:id", handlers.DeleteEmployee)

	// Customer routes
	app.Get("/customers", handlers.GetCustomers)
	app.Get("/customers/new", handlers.NewCustomer)
	app.Post("/customers", handlers.CreateCustomer)
	app.Get("/customers/:id", handlers.GetCustomer)
	app.Get("/customers/:id/edit", handlers.EditCustomer)
	app.Post("/customers/:id", handlers.UpdateCustomer)
	app.Delete("/customers/:id", handlers.DeleteCustomer)

	// Workday routes
	app.Get("/workdays/new/:employeeId", handlers.NewWorkday)
	app.Get("/workdays/:employeeId/:date", handlers.GetWorkday)
	app.Post("/workdays/:employeeId/:date", handlers.UpdateWorkday)

	// Week overview route
	app.Get("/week-overview/:employeeId", handlers.WeekOverview)

	// Daily overview route
	//app.Get("/daily-overview", handlers.DailyOverview)
	app.Get("/week", handlers.WeeklyOverview)

	// In your setupRoutes function
	app.Get("/", handleSMS)
	app.Get("/messages", showMessages)

	// Email routes
	app.Get("/email", handleIndex)
	app.Get("/emails", handleGetEmails)
	app.Post("/send", handleSendEmail)
	// In your setupRoutes function

}

func handleIndex(c *fiber.Ctx) error {
	return c.Render("index2", fiber.Map{})
}

func handleGetEmails(c *fiber.Ctx) error {
	emails, err := fetchEmails()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("email", fiber.Map{"Emails": emails})
}

func handleSendEmail(c *fiber.Ctx) error {
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

type SMS struct {
	From    string
	To      string
	Message string
}

var (
	messages []SMS
	mu       sync.Mutex
)

func handleSMS(c *fiber.Ctx) error {
	sms := SMS{
		From:    c.Query("msisdn"),
		To:      c.Query("to"),
		Message: c.Query("text"),
	}

	mu.Lock()
	messages = append(messages, sms)
	mu.Unlock()

	fmt.Printf("Received SMS:\nFrom: %s\nTo: %s\nMessage: %s\n\n", sms.From, sms.To, sms.Message)
	return c.SendStatus(fiber.StatusOK)
}

//func showMessages(c *fiber.Ctx) error {
//	mu.Lock()
//	defer mu.Unlock()
//
//	html := "<html><body><h1>SMS Messages</h1>"
//	for _, sms := range messages {
//		html += fmt.Sprintf("<p>From: %s, To: %s, Message: %s</p>", sms.From, sms.To, sms.Message)
//	}
//	html += "</body></html>"
//
//	return c.SendString(html)
//}

func showMessages(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	return c.Render("messages", fiber.Map{
		"Messages": messages,
	}, "messages")
}
