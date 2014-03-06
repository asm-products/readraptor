package mailgun

import (
	"testing"
)

// TestMessageValidity just provides some basic checks for whether a Message is valid.
// A Message needs to have a ToAddress, and then either a MimeMessage, or the set of
// FromAddress, Subject, and Body.
func TestMimeMessageValidity(t *testing.T) {
	m := MimeMessage{
		ToAddress: "bar@baz.org",
		Content:   []byte("This is my body. There are many like it but this one is mine.")}

	if m.IsValid() != true {
		t.Error("Message should have been valid!")
	}

	m.ToAddress = ""
	if m.IsValid() != false {
		t.Error("Message(2) should have been invalid!")
	}

	m = MimeMessage{ToAddress: "bar@baz.org"}
	if m.IsValid() != false {
		t.Error("Message(3) should have been invalid!")
	}
}
