package persistence

import (
	"errors"
	"github.com/lyubomirr/meme-generator-app/core/repositories"
	"gorm.io/gorm"
	"sync"
)

type uowFactory struct {}

func (u *uowFactory) Create() repositories.UnitOfWork {
	return &uow{
		db: getDB(),
	}
}

func NewUnitOfWorkFactory() repositories.UoWFactory {
	return &uowFactory{}
}

type uow struct {
	db *gorm.DB
	tx *gorm.DB
	txMux sync.Mutex
}

func (u *uow) getDb() *gorm.DB {
	u.txMux.Lock()
	defer u.txMux.Unlock()

	if u.tx != nil {
		return u.tx
	}
	return u.db
}

func (u *uow) GetUserRepository() repositories.User {
	return &mySqlUserRepository{db: u.getDb()}
}

func (u *uow) GetMemeRepository() repositories.Meme {
	return &mySqlMemeRepository{db: u.getDb()}
}

func (u *uow) GetTemplateRepository() repositories.Template {
	return &mySqlTemplateRepository{db: u.getDb()}
}

func (u *uow) GetFileRepository() repositories.File {
	return &fileRepository{}
}

func (u *uow) BeginTransaction() error {
	u.txMux.Lock()
	defer u.txMux.Unlock()

	if u.tx != nil {
		return errors.New("transaction is opened")
	}

	u.tx = u.db.Begin()
	if u.tx.Error != nil {
		return u.tx.Error
	}
	return nil
}

func (u *uow) CommitTransaction() error {
	u.txMux.Lock()
	defer u.txMux.Unlock()

	if u.tx == nil {
		return errors.New("no transaction")
	}
	db := u.tx.Commit()
	u.tx = nil
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *uow) RollbackTransaction() error {
	u.txMux.Lock()
	defer u.txMux.Unlock()
	if u.tx == nil {
		return errors.New("no transaction")
	}

	db := u.tx.Rollback()
	u.tx = nil
	if db.Error != nil {
		return db.Error
	}
	return nil
}


