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
	if user_id == "" {
		logger.Warning("Account created with no user_id")
	}
	if label == "" {
		logger.Warning("Account created with no label")
	}

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

func (a Account) Scan(server *dynamodb.Server, cio *cio_lite.ContextIOLite) {
	logger.Info("Scanning %v", a)

	folders, err := cio.GetFolders(a.UserID, a.Label, cio_lite.Params{})
	if err != nil {
		logger.Error("Unable to get folders for %v: %v", a, err)
	}

	for _, folder := range folders {
		params := cio_lite.Params{IncludeFlags: true}
		messages, err := cio.GetMessages(a.UserID, a.Label, folder.Name, params)

		if err != nil {
			logger.Error("Problem fetching messages: %v", err)
			continue
		}

		// Page through all messages in the folder
		for len(messages) > 0 {

			// Build aggregates
			senders := make(map[string]*Sender)
			for _, message := range messages {
				// Fetch sender's email address
				message_sender := message.Addresses.From.Email
				if message_sender == "" {
					logger.Warning("Unable to find sender email for message %s", message)
					continue
				}

				// Find aggregates for sender
				sender, ok := senders[message_sender]
				if !ok {
					sender = NewSender(a.UserID, a.Label, message_sender)
					senders[message_sender] = sender
				}

				// Update aggregates
				sender.TotalCount = sender.TotalCount + 1
				if message.Flags.Flags.Draft {
					sender.DraftCount = sender.DraftCount + 1
				}
				if message.Flags.Flags.Flagged {
					sender.FlaggedCount = sender.FlaggedCount + 1
				}
				if message.Flags.Flags.Answered {
					sender.AnsweredCount = sender.AnsweredCount + 1
				}
				if !message.Flags.Flags.Read {
					sender.UnreadCount = sender.UnreadCount + 1
				}
			}

			// Merge into DynamoDB
			for _, sender := range senders {	
				_, err := sender.Merge(server)

				if err != nil {
					logger.Error("Error merging sender: %v", err)
				}
			}

			// Fetch the next page of messages
			params.Offset = params.Offset + len(messages)
			messages, err = cio.GetMessages(a.UserID, a.Label, folder.Name, params)

			if strings.Contains(err.Error(), "404") {
				logger.Info("Looks like we've finished reading %v", folder)

			} else if err != nil {
				logger.Error("Problem fetching messages: %v", err)
				continue
			}
		}
	}

	logger.Info("Scan complete for %v", a)
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
