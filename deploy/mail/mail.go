package mail

import (
	"fmt"
	"ghwebhooks/config"
	"ghwebhooks/context"
	"net/smtp"
)

func createBody(context *context.Context) []byte {
	return []byte(
		fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s (%s) deployed\r\n\r\n",
			config.GetMailFrom(),
			config.GetMailTo(),
			context.AppName,
			context.Artefact.Tag))
}

func Sendmail(context *context.Context) error {
	from := config.GetMailFrom()
	to := []string{config.GetMailTo()}
	body := createBody(context)

	return smtp.SendMail(config.GetMailHost(), nil, from, to, body)
}
