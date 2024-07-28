// internal/utils/twilio.go

package utils

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

var twilioClient *twilio.RestClient

func InitTwilio(accountSid, authToken string) {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

func SendSMS(to, from, body string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)

	_, err := twilioClient.Api.CreateMessage(params)
	return err
}

func GetIncomingSMS() ([]twilioApi.ApiV2010Message, error) {
	params := &twilioApi.ListMessageParams{}
	params.SetLimit(20)

	messages, err := twilioClient.Api.ListMessage(params)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
