package database

import (
	"database/sql"

	"github.com/GabrielBrotas/eda-events/internal/entity"
)

type TransactionDB struct {
	tx *sql.Tx
}

func NewTransactionDB(tx *sql.Tx) *TransactionDB {
	return &TransactionDB{
		tx: tx,
	}
}

func (t *TransactionDB) Create(transaction *entity.Transaction) error {
	stmt, err := t.tx.Prepare("INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(transaction.ID, transaction.AccountFrom.ID, transaction.AccountTo.ID, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
