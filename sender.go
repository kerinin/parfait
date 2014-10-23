package main

import (
	"fmt"
	"strconv"
	"strings"

	"crypto/md5"

	"github.com/crowdmob/goamz/dynamodb"
)

const SenderTableName = "parfait_senders"

type Sender struct {
	Address string `json:address`
	userID  string
	label   string

	TotalCount    uint `json:total_count`
	UnreadCount   uint `json:unread_count`
	AnsweredCount uint `json:answered_count`
	FlaggedCount  uint `json:flagged_count`
	DraftCount    uint `json:draft_count`
}

func NewSender(user_id string, label string, address string) *Sender {
	if user_id == "" {
		logger.Warning("Account created with no user_id")
	}
	if label == "" {
		logger.Warning("Account created with no label")
	}
	if address == "" {
		logger.Warning("Account created with no address")
	}

	return &Sender{Address: address, userID: user_id, label: label}
}

func BootstrapSenderDynamoDB(server *dynamodb.Server) error {
	var err error

	// Bootstrap DynamoDB
	_, err = server.DescribeTable(SenderTableName)
	if err != nil && strings.Contains(err.Error(), "ResourceNotFoundException") {
		table_description := dynamodb.TableDescriptionT{
			TableName: SenderTableName,
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

func (s Sender) String() string {
	return fmt.Sprintf("<%v:%d:%d:%d:%d:%d>", s.Address, s.TotalCount, s.UnreadCount, s.AnsweredCount, s.FlaggedCount, s.DraftCount)
}

func (s Sender) Merge(server *dynamodb.Server) (bool, error) {
	logger.Info("Merging %v into dynamo", s)

	t := s.dynamoTable(server)

	attributes := []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("total_count", fmt.Sprintf("%d", s.TotalCount)),
		*dynamodb.NewNumericAttribute("unread_count", fmt.Sprintf("%d",s.UnreadCount)),
		*dynamodb.NewNumericAttribute("answered_count", fmt.Sprintf("%d", s.AnsweredCount)),
		*dynamodb.NewNumericAttribute("flagged_count", fmt.Sprintf("%d", s.FlaggedCount)),
		*dynamodb.NewNumericAttribute("draft_count", fmt.Sprintf("%d", s.DraftCount)),
	}

	return t.AddAttributes(s.dynamoKey(), attributes)
}

func (s *Sender) Load(server *dynamodb.Server) error {
	logger.Info("Loading %v from dynamo", s)

	var err error
	var i uint64

	t := s.dynamoTable(server)
	attributes, err := t.GetItem(s.dynamoKey())

	if err != nil {
		return err
	}

	i, err = strconv.ParseUint(attributes["total_count"].Value, 10, 0)
	if err != nil {
		return err
	}
	s.TotalCount = uint(i)

	i, err = strconv.ParseUint(attributes["unread_count"].Value, 10, 0)
	if err != nil {
		return err
	}
	s.UnreadCount = uint(i)

	i, err = strconv.ParseUint(attributes["answered_count"].Value, 10, 0)
	if err != nil {
		return err
	}
	s.AnsweredCount = uint(i)

	i, err = strconv.ParseUint(attributes["flagged_count"].Value, 10, 0)
	if err != nil {
		return err
	}
	s.FlaggedCount = uint(i)

	i, err = strconv.ParseUint(attributes["draft_count"].Value, 10, 0)
	if err != nil {
		return err
	}
	s.DraftCount = uint(i)

	return nil
}

func (s Sender) dynamoKey() *dynamodb.Key {
	raw_key := fmt.Sprintf("%v:%v:%v", s.Address, s.userID, s.label)
	key_bytes := md5.Sum([]byte(raw_key))
	key := fmt.Sprintf("sender:%x", key_bytes)

	return &dynamodb.Key{HashKey: key}
}

func (s Sender) dynamoTable(server *dynamodb.Server) *dynamodb.Table {
	attribute := dynamodb.NewStringAttribute("id", "")
	primary_key := dynamodb.PrimaryKey{KeyAttribute: attribute}

	return server.NewTable("parfait_senders", primary_key)
}
