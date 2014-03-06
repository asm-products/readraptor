package mailgun

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client is the base configuration for the Mailgun connection.
type Client struct {
	ApiKey   string
	Domain   string
	Hostname string // this is mostly for testing purposes
}

// MailgunMessage is used for those things that are "sendable"
// to a Mailgun endpoint
type MailgunMessage interface {
	IsValid() bool
	GetRequest(Client) *http.Request
	Endpoint() string
}

// NewClient creates a new Client using the API key and Domain passed in
func NewClient(apikey, domain string) *Client {
	return &Client{ApiKey: apikey, Domain: domain, Hostname: "https://api.mailgun.net"}
}

// Endpoint returns the URL to the Mailgun API for the MailgunMessage type that
// is being used.
func (mailgun Client) Endpoint(m MailgunMessage) string {
	return fmt.Sprintf("%s/v2/%s/%s", mailgun.Hostname, mailgun.Domain, m.Endpoint())
}

// Send processes a MailgunMessage and emits it to Mailgun for processing. If
// there is no error, it will return the text that the Mailgun endpoint returns.
func (mailgun Client) Send(message MailgunMessage) (result string, err error) {
	client := &http.Client{}

	if !message.IsValid() {
		log.Print("Mailgun.Send did not receive a valid Message object!")
		return
	}

	request := message.GetRequest(mailgun)
	request.SetBasicAuth("api", mailgun.ApiKey)
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Failed to send request: ", err)
		return
	}
	defer response.Body.Close()

	body_bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Failed to read response: ", err)
		return
	}

	return string(body_bytes), nil
}
