package main

import (
	"fmt"
	"strings"
	// "time"

	"crypto/md5"

	"github.com/crowdmob/goamz/dynamodb"
	"github.com/kerinin/parfait/cio_lite"
)

const AccountTableName = "parfait_accounts"

type Account struct {
	UserID string
	Label  string
}

// func PartiallyScannedAccounts() []Account {
// }

func NewAccount(user_id string, label string) Account {
	return Account{UserID: user_id, Label: label}
}

func BootstrapAccountDynamoDB(server *dynamodb.Server) error {
	var err error

	// Bootstrap DynamoDB
	_, err = server.DescribeTable(AccountTableName)
	if err != nil && strings.Contains(err.Error(), "ResourceNotFoundException") {
		table_description := dynamodb.TableDescriptionT{
			TableName: AccountTableName,
			ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
				ReadCapacityUnits:  4,
				WriteCapacityUnits: 20,
			},
			KeySchema: []dynamodb.KeySchemaT{
				dynamodb.KeySchemaT{AttributeName: "id", KeyType: "HASH"},
			},
			AttributeDefinitions: []dynamodb.AttributeDefinitionT{
				dynamodb.AttributeDefinitionT{Name: "id", Type: "S"},
			},
		}
		ok, err := server.CreateTable(table_description)
		if err != nil {
			return err
		}
		logger.Info("CreateTable says %v", ok)

	} else if err != nil {
		return err
	}

	return nil
}

func (a Account) String() string {
	return fmt.Sprintf("<%v:%v>", a.UserID, a.Label)
}

func (a Account) Save(server *dynamodb.Server) error {
	logger.Info("Saving %v", a)

	t := a.dynamoTable(server)

	attributes := []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("user_id", a.UserID),
		*dynamodb.NewStringAttribute("label", a.Label),
	}

	_, err := t.PutItem(a.dynamoKeyString(), "", attributes)
	if err != nil {
		return err
	}

	return nil
}

// Returns false if the account doesn't exist
func (a *Account) Load() (bool, error) {
	return true, nil
}

func (a Account) Scan(server *dynamodb.Server, cio *cio_lite.ContextIOLiteAPI) error {
	logger.Info("Scanning %v", a)

	folders, err := cio.GetFolders(a.UserID, a.Label, cio_lite.Params{})
	if err != nil {
		return err
	}

	for _, folder := range folders {
		logger.Warning("Not really scanning folder %v", folder)
	}

	/*
		for {
			senders := make(map[string]Sender)

			for message := range messages {
				// merge message into senders hash
			}

			for _, sender := range senders {
				sender.Merge(server)
			}
		}
	*/
	return nil
}

func (a Account) dynamoKeyString() string {
	raw_key := fmt.Sprintf("%v:%v", a.UserID, a.Label)
	key_bytes := md5.Sum([]byte(raw_key))

	return fmt.Sprintf("account:%x", key_bytes)
}

func (a Account) dynamoKey() dynamodb.Key {
	return dynamodb.Key{HashKey: a.dynamoKeyString()}
}

func (a Account) dynamoTable(server *dynamodb.Server) *dynamodb.Table {
	attribute := dynamodb.NewStringAttribute("id", "")
	primary_key := dynamodb.PrimaryKey{KeyAttribute: attribute}

	return server.NewTable("parfait_accounts", primary_key)
}
