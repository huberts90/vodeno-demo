package sender

import (
	"context"
	"vodeno.com/demo/internal/store"
)

type Message struct {
	MessageId int
	Email     string
	Title     string
	Content   string
	MailingId int
}

type JobSeeker interface {
	GetScheduledMailing(ctx context.Context) ([]Message, error)
	store.MessageDeleter
	store.MailingDeleter
}

type Sender interface {
	Send(context.Context, JobSeeker) error
}
