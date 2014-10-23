package main

import (
	"github.com/crowdmob/goamz/dynamodb"
	"github.com/go-martini/martini"
	"github.com/kerinin/parfait/cio_lite"
)

func RunAPI(server *dynamodb.Server, cio *cio_lite.ContextIOLite) {
	m := martini.Classic()

	m.Post("/users/:user_id/email_accounts/:label", func(params martini.Params) (int, string) {
		var err error
		account := NewAccount(params["user_id"], params["label"])

		_, err = cio.GetEmailAccount(account.UserID, account.Label, cio_lite.Params{})
		if err != nil {
			return 500, err.Error()
		}

		// Save user
		if err = account.Save(server); err != nil {
			return 500, err.Error()
		}

		// Start scanning
		go account.Scan(server, cio)

		// Respond to request
		return 200, ""
	})

	m.Run()
}
