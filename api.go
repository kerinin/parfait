package main

import (
	"github.com/go-martini/martini"
	"github.com/crowdmob/goamz/dynamodb"
)

func RunAPI(server *dynamodb.Server) {
	m := martini.Classic()

	m.Post("/users/:user_id/email_accounts/:label", func(params martini.Params) (int, string) {
		account := NewAccount(params["user_id"], params["label"])

		// Save user
		if err := account.Save(server); err != nil {
			return 500, err.Error()
		}

		// Start scanning
		go account.Scan(server)

		// Respond to request
		return 200, ""
	})

	m.Run()
}
