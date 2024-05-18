package mock

import (
	"context"
	"go.uber.org/zap"
	"vodeno.com/demo/internal/sender"
)

type mocker struct {
	logger *zap.Logger
}

func NewMocker(logger *zap.Logger) *mocker {
	return &mocker{logger: logger}
}

func (m *mocker) Send(ctx context.Context, js sender.JobSeeker) {
	m.logger.Info("Mock sender has started")

	var mailingIDs = make(map[int]int, 0)
	msgList, err := js.GetScheduledMailing(ctx)
	if err != nil {
		m.logger.Error("Failed to retrieve scheduled mailing", zap.Error(err))
	}
	// Send
	for _, msg := range msgList {
		m.logger.Info("Message has been sent", zap.String("Email", msg.Email), zap.String("Title", msg.Title))
		err := js.DeleteMessage(ctx, msg.MessageId)
		// TODO: Improve the error handling
		if err != nil {
			m.logger.Error("Failed to deleted message", zap.Int("Id", msg.MessageId), zap.Error(err))
		}
		_, ok := mailingIDs[msg.MailingId]
		if !ok {
			mailingIDs[msg.MailingId] = msg.MailingId
		}
	}

	// Delete mailing
	for id := range mailingIDs {
		err = js.DeleteMailing(ctx, id)
		if err != nil {
			m.logger.Error("Failed to deleted mailing", zap.Int("Id", id), zap.Error(err))
		}
	}

	m.logger.Info("Mock sender has finished")
}
