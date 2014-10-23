package cio_lite

type Address struct {
	Email string `json:email`
	Name  string `json:email`
}

type Attachment struct {
	Size               uint   `json:size`
	Type               string `json:type`
	FileName           string `json:file_name`
	BodySection        string `json:body_section`
	ContentDisposition string `json:content_disposition`
	EmailMessageID     string `json:email_message_id`
}

type Message struct {
	SentAt         uint                         `json:sent_at`
	ReceivedAt     uint                         `json:received_at`
	Addresses      map[string]Address           `json:addresses`
	PersonInfo     map[string]map[string]string `json:person_info`
	EmailMessageID string                       `json:email_message_id`
	Attachments    []Attachment                 `json:attachments`
	Subject        string                       `json:subject`
	Folders        []string                     `json:folders`
	EmailAccounts  []map[string]string          `json:email_accounts`
}
