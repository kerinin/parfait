package cio_lite

type User struct {
	ID              string   `json:id`
	Username        string   `json:username`
	Created         int      `json:created`
	Suspended       int      `json:suspended`
	EmailAddresses  []string `json:email_addresses`
	FirstName       string   `json:first_name`
	LastName        string   `json:last_name`
	PasswordExpired int      `json:password_expired`
	// EmailAccounts []EmailAccount `json:email_accounts`
}
