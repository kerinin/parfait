package cio_lite

import (
	"fmt"
	"regexp"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

const ctxio = "https://api.context.io"

var successRegex = regexp.MustCompile(`^2`)

type ContextIOLite struct {
	apiKey    string
	apiSecret string
}

func NewContextIOLite(key, secret string) *ContextIOLite {
	return &ContextIOLite{apiKey: key, apiSecret: secret}
}

func (cio ContextIOLite) GetUsers(params Params) ([]User, error) {
	url := fmt.Sprintf("%v/lite/users%s", ctxio, params.QueryString())
	var users []User

	if err := cio.request("GET", url, &users); err != nil {
		return users, err
	}

	return users, nil
}

func (cio ContextIOLite) GetEmailAccount(id, label string, params Params) (EmailAccount, error) {
	url := fmt.Sprintf("%v/lite/users/%s/email_accounts/%s%s", ctxio, id, url.QueryEscape(label), params.QueryString())
	var email_account EmailAccount

	if err := cio.request("GET", url, &email_account); err != nil {
		return email_account, err
	}

	return email_account, nil
}

func (cio ContextIOLite) GetFolders(id, label string, params Params) ([]Folder, error) {
	url := fmt.Sprintf("%v/lite/users/%s/email_accounts/%s/folders%s", ctxio, id, url.QueryEscape(label), params.QueryString())
	var folders []Folder

	if err := cio.request("GET", url, &folders); err != nil {
		return folders, err
	}

	return folders, nil
}

func (cio ContextIOLite) GetMessages(id string, label string, folder string, params Params) ([]Message, error) {
	url := fmt.Sprintf("%v/lite/users/%s/email_accounts/%s/folders/%s/messages%s", ctxio, id, url.QueryEscape(label), url.QueryEscape(folder), params.QueryString())
	var messages []Message

	if err := cio.request("GET", url, &messages); err != nil {
		return messages, err
	}

	return messages, nil
}

func (cio ContextIOLite) request(method, url string, ret interface{}) error {
	logger.Debug("Making %v request to %v", method, url)
	var err error

	req, err := http.NewRequest(method, url, nil)
	cio.sign_oauth(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("CIO API Returns error: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error reading CIO Response: %v", err)
	}

	if !successRegex.MatchString(resp.Status) {
		// logger.Debug("Received Body:\n%s", body)
		return fmt.Errorf("CIO API returned status %v", resp.Status)
	}

	if ret != nil {
		// logger.Debug("Received JSON:\n%s", body)
		err = json.Unmarshal(body, &ret)
		if err != nil {
			return fmt.Errorf("Error decoding CIO Response JSON: %v", err)
		}
	}

	return nil
}

func (cio ContextIOLite) sign_oauth(req *http.Request) {
	var client oauth.Client
	credentials := oauth.Credentials{cio.apiKey, cio.apiSecret}
	client.Credentials = credentials
	authHeaders := client.AuthorizationHeader(nil, "GET", req.URL, url.Values{})
	req.Header.Set("Authorization", authHeaders)
}
