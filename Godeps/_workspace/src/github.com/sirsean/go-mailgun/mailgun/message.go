package mailgun

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Message is the structure for communicating information to Mailgun.
type Message struct {
	FromName       string
	FromAddress    string
	ToAddress      string
	CCAddressList  []string
	BCCAddressList []string
	Subject        string
	Body           string
	AttachmentList []string
	InlineList     []string
}

// From concatenates the FromName and FromAddress of the message
func (m Message) From() string {
	return fmt.Sprintf("%s <%s>", m.FromName, m.FromAddress)
}

// BCCAddresses joins the BCCAddressList together with commas
func (m Message) BCCAddresses() string {
	return strings.Join(m.BCCAddressList, ", ")
}

// CCAddresses joins the CCAddressList together with commas
func (m Message) CCAddresses() string {
	return strings.Join(m.CCAddressList, ", ")
}

// IsValid verifies that the Message has all of the required
// fields filled in
func (message Message) IsValid() (validity bool) {
	if message.ToAddress == "" ||
		message.FromAddress == "" ||
		message.Subject == "" ||
		message.Body == "" {
		return false
	}

	return true
}

// URLValues converts the Message to a format that can be used
// for POSTing to Mailgun
func (m Message) URLValues() url.Values {
	values := make(url.Values)
	values.Set("to", m.ToAddress)
	values.Set("from", m.From())
	values.Set("subject", m.Subject)
	values.Set("text", m.Body)

	if m.CCAddresses() != "" {
		values.Set("cc", m.CCAddresses())
	}

	if m.BCCAddresses() != "" {
		values.Set("bcc", m.BCCAddresses())
	}

	for _, attachment := range m.AttachmentList {
		values.Add("attachment", attachment)
	}

	for _, inline := range m.InlineList {
		values.Add("inline", inline)
	}

	return values
}

// GetRequest returns a skeleton http.Request refernce with the Content-Type
// header filled in, along with the formatting required for this type of Message
func (message Message) GetRequest(mailgun Client) (request *http.Request) {
	request, _ = http.NewRequest("POST", mailgun.Endpoint(message), strings.NewReader(message.URLValues().Encode()))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	return
}

// Endpoint returns the final part of the path required for creating
// the Mailgun URL for this type of Message
func (message Message) Endpoint() string {
	return "messages"
}
