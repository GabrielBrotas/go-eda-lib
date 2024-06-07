package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/pkg/uow"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite

	db          *sql.DB
	uow         *uow.Uow
	client      *entity.Client
	client2     *entity.Client
	accountFrom *entity.Account
	accountTo   *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	_, err = db.Exec("CREATE TABLE clients (id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at DATE)")
	s.Nil(err)
	_, err = db.Exec("CREATE TABLE accounts (id VARCHAR(255), client_id VARCHAR(255), balance INT, created_at DATE)")
	s.Nil(err)
	_, err = db.Exec("CREATE TABLE transactions (id VARCHAR(255), account_id_from VARCHAR(255), account_id_to VARCHAR(255), amount INT, created_at DATE)")
	s.Nil(err)

	client, err := entity.NewClient("John", "j@j.com")
	s.Nil(err)
	s.client = client
	client2, err := entity.NewClient("John2", "jj@j.com")
	s.Nil(err)
	s.client2 = client2

	// Creating accounts
	accountFrom := entity.NewAccount(s.client)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(s.client2)
	accountTo.Balance = 1000
	s.accountTo = accountTo

	s.uow = uow.NewUow(db)
	s.uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return NewTransactionDB(tx)
	})
	s.uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return NewAccountDB(tx)
	})
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	_, err := s.db.Exec("DROP TABLE clients")
	s.Nil(err)
	_, err = s.db.Exec("DROP TABLE accounts")
	s.Nil(err)
	_, err = s.db.Exec("DROP TABLE transactions")
	s.Nil(err)
}
func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	err := s.uow.Do(context.Background(), func(u *uow.Uow) error {
		transactionRepo, err := u.GetRepository(context.Background(), "TransactionDB")
		if err != nil {
			return err
		}
		accountRepo, err := u.GetRepository(context.Background(), "AccountDB")
		if err != nil {
			return err
		}

		// Save accounts first
		err = accountRepo.(*AccountDB).Save(s.accountFrom)
		if err != nil {
			return err
		}
		err = accountRepo.(*AccountDB).Save(s.accountTo)
		if err != nil {
			return err
		}

		// Create transaction
		transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
		if err != nil {
			return err
		}
		return transactionRepo.(*TransactionDB).Create(transaction)
	})
	s.Nil(err)
}
