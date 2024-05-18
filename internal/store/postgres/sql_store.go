package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"vodeno.com/demo/internal/sender"
	"vodeno.com/demo/internal/store"
)

type sqlStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *sqlStore {
	return &sqlStore{
		db: db,
	}
}

// TODO: SQL injection
func (s *sqlStore) InsertMessage(ctx context.Context, message *store.Message) error {
	err := s.createMailingIfNotExists(ctx, message.MailingId)
	if err != nil {
		return fmt.Errorf("failed to retrieve a mailing: %w", err)
	}
	customerId, err := s.createCustomerIfNotExists(ctx, message.Email)
	if err != nil {
		return fmt.Errorf("failed to retrieve a customer: %w", err)
	}

	err = s.createMessage(ctx, customerId, message)
	if err != nil {
		return fmt.Errorf("failed to create a message: %w", err)
	}

	return nil
}

func (s *sqlStore) DeleteMessage(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM MESSAGES WHERE ID=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *sqlStore) OrderMailing(ctx context.Context, mj *store.MailingJob) error {
	_, err := s.db.ExecContext(ctx, "UPDATE MAILINGS SET IS_SCHEDULED=TRUE WHERE ID=$1", mj.MailingId)
	if err != nil {

		return fmt.Errorf("failed to order a mailing: %w", err)
	}

	return nil
}

func (s *sqlStore) DeleteMailing(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM MAILINGS WHERE ID=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete mailing: %w", err)
	}

	return nil
}

func (s *sqlStore) GetScheduledMailing(ctx context.Context) ([]sender.Message, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT MESSAGES.ID, C.EMAIL, TITLE, CONTENT, ML.ID as MAILING_ID "+
		"FROM MESSAGES JOIN MAILINGS ML on ML.ID = MESSAGES.MAILING_ID JOIN CUSTOMERS C on C.ID = MESSAGES.CUSTOMER_ID "+
		"WHERE ML.IS_SCHEDULED=TRUE")

	if err != nil {
		return nil, fmt.Errorf("failed to select scheduled mailing: %w", err)
	}
	defer rows.Close()

	// TODO: Get the size of list
	var msgList []sender.Message
	for rows.Next() {
		var id int
		var email string
		var title string
		var content string
		var mailingId int
		err = rows.Scan(&id, &email, &title, &content, &mailingId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve message details: %w", err)
		}
		msg := sender.Message{
			MessageId: id,
			Email:     email,
			Title:     title,
			Content:   content,
			MailingId: mailingId,
		}
		// TODO: Find more efficient way to deal with slice growing
		msgList = append(msgList, msg)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over messages: %w", err)
	}

	return msgList, err
}

func (s *sqlStore) createMailingIfNotExists(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO MAILINGS(id) VALUES($1) ON CONFLICT DO NOTHING`, id)
	return err
}

func (s *sqlStore) createCustomerIfNotExists(ctx context.Context, email string) (int, error) {
	var id int
	_, err := s.db.ExecContext(ctx, `INSERT INTO CUSTOMERS(email) VALUES($1) ON CONFLICT DO NOTHING`, email)
	if err != nil {
		return 0, err
	}
	// Instead of playing around `RETURNING ID`
	err = s.db.QueryRowContext(ctx, "SELECT id FROM CUSTOMERS WHERE email=$1", email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *sqlStore) createMessage(ctx context.Context, customerId int, m *store.Message) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO MESSAGES(mailing_id, customer_id, title, content) VALUES($1, $2, $3, $4)`,
		m.MailingId, customerId, m.Title, m.Content)

	return err
}
