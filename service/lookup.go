package service

import (
	"ghwebhooks/context"
	"ghwebhooks/types"
	"os/user"
	"strconv"
)

type UserLookup struct {
	uid, gid int
}

func Lookup(context *context.Context, status *types.Status) (lookup UserLookup) {
	username := context.AppName
	if user, err := user.Lookup(username); err != nil {
		status.Fail(err)
	} else {
		uid, err := strconv.Atoi(user.Uid)
		gid, err := strconv.Atoi(user.Gid)

		if err != nil {
			status.Fail(err)
		} else {
			status.LogF("found uid: %v, gid: %v for user: %s", uid, gid, username)
			lookup = UserLookup{
				uid,
				gid,
			}
		}
	}
	return lookup
}
