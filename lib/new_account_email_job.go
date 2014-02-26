package readraptor

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

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

	token := genKey("confirm" + account.Email)
	account.ConfirmationToken = &token

	message, err := j.CreateMessage(account)
	if err != nil {
		return err
	}

	_, err = dbmap.Exec(
		"update accounts set confirmation_token = $1, confirmation_sent_at = $2 where id = $3",
		token,
		time.Now(),
		j.AccountId,
	)
	if err != nil {
		return err
	}

	if os.Getenv("MAILGUN_API_KEY") != "" {
		mg := mailgun.NewClient(os.Getenv("MAILGUN_API_KEY"), os.Getenv("MAILGUN_DOMAIN"))
		_, err := mg.Send(message)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(message)
	}

	return nil
}

func (j *NewAccountEmailJob) CreateMessage(account *Account) (*mailgun.Message, error) {
	template, err := template.ParseFiles(os.Getenv("RR_ROOT") + "/templates/new_account_email.tmpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
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
