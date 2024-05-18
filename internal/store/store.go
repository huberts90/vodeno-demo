package store

import (
	"context"
)

type Message struct {
	Email     string
	Title     string
	Content   string
	MailingId int `json:"mailing_id"`
}

type MailingJob struct {
	MailingId int `json:"mailing_id"`
}

type MessageDeleter interface {
	DeleteMessage(context.Context, int) error
}

type MailingDeleter interface {
	DeleteMailing(context.Context, int) error
}

type Deleter interface {
	MessageDeleter
	MailingDeleter
}

type MessengerStore interface {
	InsertMessage(context.Context, *Message) error
	MessageDeleter
	OrderMailing(context.Context, *MailingJob) error
}
