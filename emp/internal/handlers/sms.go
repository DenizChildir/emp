// internal/handlers/sms.go

package handlers

import (
	"emp/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SMSPage(c *fiber.Ctx) error {
	messages, err := utils.GetIncomingSMS()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Render("sms", fiber.Map{
		"Messages": messages,
	})
}

func SendSMS(c *fiber.Ctx) error {
	to := c.FormValue("to")
	body := c.FormValue("body")
	from := "YOUR_TWILIO_PHONE_NUMBER" // Replace with your Twilio phone number

	err := utils.SendSMS(to, from, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/sms")
}

func ReceiveSMS(c *fiber.Ctx) error {
	// This endpoint will be called by Twilio when you receive an SMS
	// You can process the incoming SMS here and store it in your database if needed
	//from := c.FormValue("From")
	//body := c.FormValue("Body")

	// For now, we'll just log the incoming SMS
	//c.App().Logger().Infof("Received SMS from %s: %s", from, body)

	// Respond to Twilio with an empty TwiML response
	return c.SendString("<Response></Response>")
}
