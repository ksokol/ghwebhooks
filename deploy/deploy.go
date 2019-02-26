package deploy

import (
	"ghwebhooks/context"
	"ghwebhooks/github"
	"ghwebhooks/mail"
	"ghwebhooks/service"
	"ghwebhooks/types"
)

func Deploy(context *context.Context, status *types.Status) {
	if service.Update(context, status); status.Success != true {
		return
	}

	if mail.Sendmail(context, status); status.Success != true {
		return
	}

	if github.RemovePreviousReleases(&context.Event, status); status.Success != true {
		return
	}

	github.RemoveDraftReleases(&context.Event, status)
}
