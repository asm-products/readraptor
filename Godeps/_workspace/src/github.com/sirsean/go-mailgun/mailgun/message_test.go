package mailgun

import (
	"fmt"
	"testing"
)

func TestMessageCCAddresses(t *testing.T) {
	m := Message{}
	m.CCAddressList = []string{"john@doe.org", "jane@smith.org"}
	if m.CCAddresses() != "john@doe.org, jane@smith.org" {
		t.Error(fmt.Sprintf("m.CCAddresses returned the wrong result: '%s'", m.CCAddresses()))
	}
}

func TestMessageBCCAddresses(t *testing.T) {
	m := Message{}
	m.BCCAddressList = []string{"jane@doe.org", "john@smith.org"}
	if m.BCCAddresses() != "jane@doe.org, john@smith.org" {
		t.Error(fmt.Sprintf("m.BCCAddresses returned the wrong result: '%s'", m.BCCAddresses()))
	}
}

func TestURLValues(t *testing.T) {
	m := Message{
		FromName:       "Testy McTestsalot",
		FromAddress:    "test@foo.org",
		ToAddress:      "bar@baz.org",
		CCAddressList:  []string{"cc1@foo.org", "cc2@baz.org"},
		BCCAddressList: []string{"bcc1@baz.org", "bcc2@foo.org"},
		Subject:        "Best subject evar!",
		Body:           "This is my body. There are many like it but this one is mine.",
		AttachmentList: []string{"TWFpbGd1biB0ZXN0"},
		InlineList:     []string{"aW1hZ2UgZ29lcyBoZXJl"}}
	urlValues := m.URLValues()

	if urlValues.Get("from") != "Testy McTestsalot <test@foo.org>" {
		t.Error("'from' value not set correctly!")
	}

	if urlValues.Get("to") != "bar@baz.org" {
		t.Error("'to' value not set correctly!")
	}

	if urlValues.Get("subject") != "Best subject evar!" {
		t.Error("'subject' value not set correctly!")
	}

	if urlValues.Get("text") != "This is my body. There are many like it but this one is mine." {
		t.Error("'text' value not set correctly!")
	}

	if urlValues.Get("cc") != "cc1@foo.org, cc2@baz.org" {
		t.Error("'cc' value not set correctly!")
	}

	if urlValues.Get("bcc") != "bcc1@baz.org, bcc2@foo.org" {
		t.Error("'bcc' value not set correctly!")
	}

	if stringSlicesDifferent(urlValues["attachment"], []string{"TWFpbGd1biB0ZXN0"}) {
		t.Error("'attachment' value not set correctly!")
	}

	if stringSlicesDifferent(urlValues["inline"], []string{"aW1hZ2UgZ29lcyBoZXJl"}) {
		t.Error("'inline' value not set correctly!")
	}
}

func TestMinimalURLValues(t *testing.T) {
	m := Message{
		FromName:    "Testy McTestsalot",
		FromAddress: "test@foo.org",
		ToAddress:   "bar@baz.org",
		Subject:     "Best subject evar!",
		Body:        "This is my body. There are many like it but this one is mine."}
	urlValues := m.URLValues()

	if urlValues.Get("from") != "Testy McTestsalot <test@foo.org>" {
		t.Error("'from' value not set correctly!")
	}

	if urlValues.Get("to") != "bar@baz.org" {
		t.Error("'to' value not set correctly!")
	}

	if urlValues.Get("subject") != "Best subject evar!" {
		t.Error("'subject' value not set correctly!")
	}

	if urlValues.Get("text") != "This is my body. There are many like it but this one is mine." {
		t.Error("'text' value not set correctly!")
	}

	if urlValues.Get("cc") != "" {
		t.Error("'cc' value not set correctly!")
	}

	if urlValues.Get("bcc") != "" {
		t.Error("'bcc' value not set correctly!")
	}

	if urlValues.Get("attachment") != "" {
		t.Error("'attachment' value not set correctly!")
	}

	if urlValues.Get("inline") != "" {
		t.Error("'inline' value not set correctly!")
	}
}

// TestMessageValidity just provides some basic checks for whether a Message is valid.
// A Message needs to have a ToAddress, FromAddress, Subject, and Body.
func TestMessageValidity(t *testing.T) {
	m := Message{
		FromName:    "Testy McTestsalot",
		FromAddress: "test@foo.org",
		ToAddress:   "bar@baz.org",
		Subject:     "Best subject evar!",
		Body:        "This is my body. There are many like it but this one is mine."}

	if m.IsValid() != true {
		t.Error("Message should have been valid!")
	}

	m.ToAddress = ""
	if m.IsValid() != false {
		t.Error("Message(2) should have been invalid!")
	}

	m = Message{ToAddress: "bar@baz.org"}
	if m.IsValid() != false {
		t.Error("Message(3) should have been invalid!")
	}
}

func stringSlicesDifferent(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return true
	}

	for index, value := range slice1 {
		if value != slice2[index] {
			return true
		}
	}

	return false
}
