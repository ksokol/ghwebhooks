package mail

import (
	"fmt"
	"ghwebhooks/config"
	"ghwebhooks/types"
	"net/smtp"
)

func createBody(context *types.Context) []byte {
	return []byte(
		fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s (%s) deployed\r\n\r\n",
			config.GetMailFrom(),
			config.GetMailTo(),
			context.AppName,
			context.Artefact.Tag))
}

func Sendmail(context *types.Context) error {
	from := config.GetMailFrom()
	to := []string{config.GetMailTo()}
	body := createBody(context)

	return smtp.SendMail(config.GetMailHost(), nil, from, to, body)
}
