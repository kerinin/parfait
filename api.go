package main

import (
	"github.com/crowdmob/goamz/dynamodb"
	"github.com/go-martini/martini"
	"github.com/kerinin/parfait/cio_lite"
)

func RunAPI(server *dynamodb.Server, cio *cio_lite.ContextIOLiteAPI) {
	m := martini.Classic()

	m.Post("/users/:user_id/email_accounts/:label", func(params martini.Params) (int, string) {
		var err error
		account := NewAccount(params["user_id"], params["label"])

		// NOTE: Might want to check CIO user exists before saving

		// Save user
		if err = account.Save(server); err != nil {
			return 500, err.Error()
		}

		// Start scanning
		if err = account.Scan(server, cio); err != nil {
			return 500, err.Error()
		}

		// Respond to request
		return 200, ""
	})

	m.Run()
}
