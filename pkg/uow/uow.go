// Unit Of Work
// It helps maintain data consistency by ensuring that all operations within a transaction are completed successfully or none are completed at all.
package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// RepositoryFactory is a function that returns a repository instance with a transaction
type RepositoryFactory func(tx *sql.Tx) interface{}

// Interface that defines the methods that a UoW must implement
type UowInterface interface {
	Register(name string, fc RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(uow *Uow) error) error
	CommitOrRollback() error
	Rollback() error
	UnRegister(name string)
}

// main struct that implements the UoW interface
type Uow struct {
	db           *sql.DB                      // Database connection
	tx           *sql.Tx                      // Transaction
	repositories map[string]RepositoryFactory // Map of repositories

}

func NewUow(db *sql.DB) *Uow {
	return &Uow{
		db:           db,
		repositories: make(map[string]RepositoryFactory),
	}
}

// Register a repository factory with a given name
func (u *Uow) Register(name string, fc RepositoryFactory) {
	u.repositories[name] = fc
}

// Unregister unregisters a repository factory with a given name
func (u *Uow) UnRegister(name string) {
	delete(u.repositories, name)
}

// GetRepository retrieves a repository by name and returns an instance of it with a transaction
func (u *Uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
	factory, exists := u.repositories[name]

	if !exists {
		return nil, fmt.Errorf("repository %s not found", name)
	}

	// Start transaction if not already started
	if u.tx == nil {
		tx, err := u.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.tx = tx
	}

	// Create a new repository with the current transaction
	repo := factory(u.tx)
	return repo, nil
}

// Do method starts a transaction and executes the function passed as a parameter. If an error occurs, it rolls back the transaction.
func (u *Uow) Do(ctx context.Context, fn func(Uow *Uow) error) error {
	// Check if a transaction has already been started
	if u.tx != nil {
		// if a transaction has already been started, it means that the Do method has already been called
		return fmt.Errorf("transaction already started")
	}

	// Start transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx

	// Execute the provided function
	err = fn(u)

	// Check if an error occurred
	if err != nil {
		// If an error occurred, rollback the transaction
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}

	// Commit the transaction
	return u.CommitOrRollback()
}

// CommitOrRollback method commits the transaction. If an error occurs, it rolls back the transaction.
func (u *Uow) CommitOrRollback() error {
	if u.tx == nil {
		return errors.New("no transaction to commit or rollback")
	}

	if err := u.tx.Commit(); err != nil {
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	u.tx = nil
	return nil
}

// Rollback rolls back the transaction
func (u *Uow) Rollback() error {
	if u.tx == nil {
		return errors.New("no transaction to rollback")
	}

	if err := u.tx.Rollback(); err != nil {
		return err
	}

	u.tx = nil
	return nil
}
