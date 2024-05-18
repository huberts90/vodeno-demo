package postgres

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *sqlStore) Insert(ctx context.Context, message *store.Message) error {
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
