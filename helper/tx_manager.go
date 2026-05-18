package helper

import "gorm.io/gorm"

type Tx interface {
	Rollback()
	Commit() error
	CheckPanic()
	Get() *gorm.DB
}

type tx struct {
	db *gorm.DB
}

func (t *tx) Rollback() {
	t.db.Rollback()
}

func (t *tx) Commit() error {
	return t.db.Commit().Error
}

func (t *tx) Get() *gorm.DB {
	return t.db
}

func (t *tx) CheckPanic() {
	if r := recover(); r != nil {
		t.Rollback()
		panic(r)
	}
}

type TxManager interface {
	New() Tx
}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManager {
	return &txManager{db: db}
}

func (txm *txManager) New() Tx {
	transaction := txm.db.Begin()
	return &tx{db: transaction}
}
