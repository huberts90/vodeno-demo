package store

import "context"

type Message struct {
	Email     string
	Title     string
	Content   string
	MailingId int `json:"mailing_id"`
}

type MailingJob struct {
	MailingId int
}

type MessengerStore interface {
	Insert(context.Context, *Message) error
	//Order(context.Context, *MailingJob) error
}
