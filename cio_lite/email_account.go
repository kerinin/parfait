package cio_lite

type EmailAccount struct {
	Server             string `json:server`
	Label              string `json:label`
	Username           string `json:username`
	Port               uint   `json:port`
	AuthenticationType string `json:authentication_type`
	Status             string `json:status`
	UseSSL             bool   `json:use_ssl`
	Type               string `json:type`
}
