package cio_lite

type Folder struct {
	Name               string `json:name`
	SymbolicName       string `json:symbolic_name`
	MessageCount       uint   `json:nb_messages`
	UnreadMessageCount uint   `json:nb_unseen_messages`
	Delimiter          string `json:delimiter`
}
