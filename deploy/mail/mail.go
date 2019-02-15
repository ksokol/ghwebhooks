package mail

import (
	"fmt"
	"ghwebhooks/types"
	"net/smtp"
)

func createBody(context *types.Context) []byte {
	return []byte(
		fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s (%s) deployed\r\n\r\n",
			context.Mail.From,
			context.Mail.To,
			context.AppName,
			context.Artefact.Tag))
}

func Sendmail(context *types.Context) error {
	from := context.Mail.From
	to := []string{context.Mail.To}
	body := createBody(context)

	return smtp.SendMail(context.Mail.Host, nil, from, to, body)
}
