package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/GabrielBrotas/eda-events/internal/entity"
	"github.com/GabrielBrotas/eda-events/pkg/uow"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite

	db     *sql.DB
	uow    *uow.Uow
	client *entity.Client
}

func (s *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	_, err = db.Exec("CREATE TABLE clients (id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at DATE)")
	s.Nil(err)
	_, err = db.Exec("CREATE TABLE accounts (id VARCHAR(255), client_id VARCHAR(255), balance INT, created_at DATE)")
	s.Nil(err)

	s.uow = uow.NewUow(db)
	s.uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return NewAccountDB(tx)
	})

	s.client, _ = entity.NewClient("John", "j@j.com")
}

func (s *AccountDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	_, err := s.db.Exec("DROP TABLE clients")
	s.Nil(err)
	_, err = s.db.Exec("DROP TABLE accounts")
	s.Nil(err)
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	err := s.uow.Do(context.Background(), func(u *uow.Uow) error {
		accountRepo, err := u.GetRepository(context.Background(), "AccountDB")
		if err != nil {
			return err
		}

		account := entity.NewAccount(s.client)
		return accountRepo.(*AccountDB).Save(account)
	})
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	_, err := s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt)
	s.Nil(err)

	err = s.uow.Do(context.Background(), func(u *uow.Uow) error {
		accountRepo, err := u.GetRepository(context.Background(), "AccountDB")
		if err != nil {
			return err
		}
		account := entity.NewAccount(s.client)
		err = accountRepo.(*AccountDB).Save(account)
		if err != nil {
			return err
		}
		accountDB, err := accountRepo.(*AccountDB).FindByID(account.ID)
		if err != nil {
			return err
		}
		s.Equal(account.ID, accountDB.ID)
		s.Equal(account.Client.ID, accountDB.Client.ID)
		s.Equal(account.Balance, accountDB.Balance)
		s.Equal(account.Client.ID, accountDB.Client.ID)
		s.Equal(account.Client.Name, accountDB.Client.Name)
		s.Equal(account.Client.Email, accountDB.Client.Email)
		return nil
	})
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestUpdateBalance() {
	_, err := s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt)
	s.Nil(err)

	err = s.uow.Do(context.Background(), func(u *uow.Uow) error {
		accountRepo, err := u.GetRepository(context.Background(), "AccountDB")
		if err != nil {
			return err
		}
		// 1 - Create a new account
		account := entity.NewAccount(s.client)
		err = accountRepo.(*AccountDB).Save(account)
		if err != nil {
			return err
		}

		// 2 - Update the account balance
		account.Credit(100)
		err = accountRepo.(*AccountDB).UpdateBalance(account)
		if err != nil {
			return err
		}

		// 3 - Check if the balance was updated
		accountDB, err := accountRepo.(*AccountDB).FindByID(account.ID)
		if err != nil {
			return err
		}
		s.Equal(float64(100), accountDB.Balance)
		return nil
	})
	s.Nil(err)
}
