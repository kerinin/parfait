package main

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
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
	}

	if err = BootstrapAccountDynamoDB(server); err != nil {
		logger.Error("Error bootstrapping Account DynamoDB: %v", err)
	}

	// Resume scanning
	// for _, account := range PartiallyScannedAccounts() {
	// 	go account.Scan(server)
	// }

	// Serve HTTP Requests
	RunAPI(server)
}
