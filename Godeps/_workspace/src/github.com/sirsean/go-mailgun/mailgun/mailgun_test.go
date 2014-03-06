package mailgun

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var mailgun = NewClient("fake API key", "testdomain.org")

func TestSendWithGoodServer(t *testing.T) {
	message := Message{
		FromName:       "Testy McTestsalot",
		FromAddress:    "test@foo.org",
		ToAddress:      "bar@baz.org",
		CCAddressList:  []string{"cc1@foo.org", "cc2@baz.org"},
		BCCAddressList: []string{"bcc1@baz.org", "bcc2@foo.org"},
		Subject:        "Best subject evar!",
		Body:           "This is my body. There are many like it but this one is mine.",
		AttachmentList: []string{"TWFpbGd1biB0ZXN0"},
		InlineList:     []string{"aW1hZ2UgZ29lcyBoZXJl"}}

	test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatal("Error when parsing the form: ", err)
		}

		recvToExpected := map[string]string{
			r.FormValue("from"):       message.From(),
			r.FormValue("to"):         message.ToAddress,
			r.FormValue("cc"):         message.CCAddresses(),
			r.FormValue("bcc"):        message.BCCAddresses(),
			r.FormValue("subject"):    message.Subject,
			r.FormValue("text"):       message.Body,
			r.FormValue("attachment"): strings.Join(message.AttachmentList, ","),
			r.FormValue("inline"):     strings.Join(message.InlineList, ",")}

		for received, expected := range recvToExpected {
			if received != expected {
				t.Error("Received '" + received + "' instead of '" + expected + "'!")
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Queued. Thank you.", "id":"fake_id_here@mailgun.org"}`)
	}))
	defer test_server.Close()

	mailgun.Hostname = test_server.URL
	mailgun.Send(message)
}

func TestMimeSendWithGoodServer(t *testing.T) {
	mime_message := "test_data/message.mime"
	file, err := os.Open(mime_message)
	if err != nil {
		t.Fatal("Error when opening 'testdata/message.mime'")
	}

	fileinfo, err := file.Stat()
	if err != nil {
		t.Fatal("Error when stat'ing \"testdata/message.mime\"")
	}

	mimeContent := make([]byte, fileinfo.Size())
	_, err = file.Read(mimeContent)
	if err != nil {
		t.Fatal("Could not read from " + mime_message)
	}

	message := MimeMessage{
		ToAddress: "to_address@bar.org",
		Content:   mimeContent}

	test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("to") != message.ToAddress {
			t.Error("Received '" + r.FormValue("to") + "' instead of '" + message.ToAddress + "'!")
		}

		file, _, _ := r.FormFile("message")
		fileContents, _ := ioutil.ReadAll(file)
		if string(fileContents) != string(message.Content) {
			t.Error("Received '\n" + string(fileContents) + "' instead of '" + string(message.Content) + "'!")
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Queued. Thank you.", "id":"fake_id_here@mailgun.org"}`)
	}))
	defer test_server.Close()

	mailgun.Hostname = test_server.URL
	mailgun.Send(message)
}
