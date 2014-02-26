package readraptor

import (
	"bytes"
	"os"
	"text/template"

	"github.com/sirsean/go-mailgun/mailgun"
)

type NewAccountEmailJob struct {
	AccountId int64
}

func (j *NewAccountEmailJob) Perform() error {
	account, err := FindAccount(j.AccountId)
	if err != nil {
		return err
	}

	message, err := j.CreateMessage(account)
	if err != nil {
		return err
	}

	if os.Getenv("MAILGUN_API_KEY") != "" {
		mg := mailgun.NewClient(os.Getenv("MAILGUN_API_KEY"), os.Getenv("MAILGUN_DOMAIN"))
		_, err := mg.Send(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *NewAccountEmailJob) CreateMessage(account *Account) (*mailgun.Message, error) {
	var buf bytes.Buffer
	template, err := template.ParseFiles("../templates/new_account_email.tmpl")
	if err != nil {
		return nil, err
	}

	err = template.Execute(&buf, account)
	if err != nil {
		return nil, err
	}

	message := &mailgun.Message{
		FromName:    "Read Raptor",
		FromAddress: "rawr@readraptor.com",
		ToAddress:   account.Email,
		Subject:     "Get started with Read Raptor",
		Body:        buf.String(),
	}

	return message, nil
}
