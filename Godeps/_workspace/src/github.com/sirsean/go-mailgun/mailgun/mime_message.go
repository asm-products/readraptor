package mailgun

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// MimeMessage is the structure for communicating a MIME message to Mailgun.
type MimeMessage struct {
	ToAddress string
	Content   []byte
}

// IsValid verifies that the Message has all of the required
// fields filled in
func (message MimeMessage) IsValid() (validity bool) {
	if message.ToAddress == "" || string(message.Content) == "" {
		return false
	}

	return true
}

// MimeReader returns a reader from which the MIME email may be read. Mailgun
// requires a different header and multipart message when talking to the MIME
// endpoint.
func (message MimeMessage) MimeReader() (b io.Reader, boundary string) {
	buffer := new(bytes.Buffer)
	mimeWriter := multipart.NewWriter(buffer)
	boundary = mimeWriter.Boundary()

	go func() {
		defer mimeWriter.Close()
		mimeWriter.WriteField("to", message.ToAddress)

		messageField, err := mimeWriter.CreateFormFile("message", "message.mime")
		if err != nil {
			log.Fatal("Could not create MIME part for the 'message' field!")
		}
		messageField.Write(message.Content)
	}()

	return buffer, boundary
}

// GetRequest returns a skeleton http.Request refernce with the Content-Type
// header filled in, along with the formatting required for this type of Message
func (message MimeMessage) GetRequest(mailgun Client) (request *http.Request) {
	mimeReader, boundary := message.MimeReader()
	request, _ = http.NewRequest("POST", mailgun.Endpoint(message), mimeReader)
	request.Header.Set("content-type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	return
}

// Endpoint returns the final part of the path required for creating
// the Mailgun URL for this type of Message
func (message MimeMessage) Endpoint() string {
	return "messages.mime"
}
