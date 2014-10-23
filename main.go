package main

import (
	"os"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
	"github.com/kerinin/parfait/cio_lite"
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		logger.Error("Credentials error: %v", err)
		return
	}
	region := aws.GetRegion("us-east-1")
	server := dynamodb.New(auth, region)

	if err = BootstrapSenderDynamoDB(server); err != nil {
		logger.Error("Error bootstrapping Sender DynamoDB: %v", err)
		return
	}

	if err = BootstrapAccountDynamoDB(server); err != nil {
		logger.Error("Error bootstrapping Account DynamoDB: %v", err)
		return
	}

	cio_key := os.Getenv("CIO_DEVELOPER_KEY")
	cio_secret := os.Getenv("CIO_DEVELOPER_SECRET")
	if cio_key == "" || cio_secret == "" {
		logger.Error("CIO_DEVELOPER_KEY and CIO_DEVELOPER_SECRET must be set in your environment (export CIO_DEVELOPER_KEY=foo)")
		return
	}
	cio := cio_lite.NewContextIOLite(cio_key, cio_secret)
	logger.Debug("Checking CIO credentials...")
	users, err := cio.GetUsers(cio_lite.Params{})
	logger.Debug("Found %d users", len(users))
	if err != nil {
		logger.Error("Problem connecting to CIO: %v", err)
		return
	}

	// Resume scanning
	// for _, account := range PartiallyScannedAccounts() {
	// 	go account.Scan(server)
	// }

	// Serve HTTP Requests
	RunAPI(server, cio)
}
