package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	sqlStore := NewSQLStore(db)
	require.NotNil(t, sqlStore)

	ctx := context.TODO()
	id := 1

	// Delete fails
	mock.ExpectExec(`DELETE FROM MESSAGES WHERE ID=\$1`).WithArgs(id).WillReturnError(errors.New("row not found"))
	err = sqlStore.DeleteMessage(ctx, id)
	require.NotNil(t, err)

	// Delete succeeds
	mock.ExpectExec(`DELETE FROM MESSAGES WHERE ID=\$1`).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	err = sqlStore.DeleteMessage(ctx, id)
	require.Nil(t, err)
}
