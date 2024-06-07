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

type ClientDBTestSuite struct {
	suite.Suite
	db  *sql.DB
	uow *uow.Uow
}

func (s *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	_, err = db.Exec("CREATE TABLE clients (id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at DATE)")
	s.Nil(err)

	s.uow = uow.NewUow(db)
	s.uow.Register("ClientDB", func(tx *sql.Tx) interface{} {
		return NewClientDB(tx)
	})
}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	_, err := s.db.Exec("DROP TABLE clients")
	s.Nil(err)
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}
func (s *ClientDBTestSuite) TestSave() {
	err := s.uow.Do(context.Background(), func(u *uow.Uow) error {
		clientRepo, err := u.GetRepository(context.Background(), "ClientDB")
		if err != nil {
			return err
		}
		client := &entity.Client{
			ID:    "1",
			Name:  "Test",
			Email: "j@j.com",
		}
		return clientRepo.(*ClientDB).Save(client)
	})
	s.Nil(err)
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John", "j@j.com")
	err := s.uow.Do(context.Background(), func(u *uow.Uow) error {
		clientRepo, err := u.GetRepository(context.Background(), "ClientDB")
		if err != nil {
			return err
		}
		return clientRepo.(*ClientDB).Save(client)
	})
	s.Nil(err)

	err = s.uow.Do(context.Background(), func(u *uow.Uow) error {
		clientRepo, err := u.GetRepository(context.Background(), "ClientDB")
		if err != nil {
			return err
		}
		clientDB, err := clientRepo.(*ClientDB).Get(client.ID)
		if err != nil {
			return err
		}
		s.Equal(client.ID, clientDB.ID)
		s.Equal(client.Name, clientDB.Name)
		s.Equal(client.Email, clientDB.Email)
		return nil
	})
	s.Nil(err)
}
