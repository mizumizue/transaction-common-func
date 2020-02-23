package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDBConn() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", "mysql:password@/test?charset=utf8&parseTime=True&loc=Local")
}

func Transact(ctx context.Context, db *sqlx.DB, txFunc func(*sqlx.Tx) error) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		err = BeginTxErr{err: err}
		return
	}
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				panic(fmt.Errorf("rollback failed. detail: %v", err))
			}
			panic(p) // re-throw panic after Rollback
		}
		if err != nil && errors.As(err, &BeginTxErr{}) {
			return
		}

		if err != nil {
			txErr := tx.Rollback() // err is nil; if Commit returns error update err
			if txErr != nil {
				err = fmt.Errorf("transaction failed. detail: %w", err)
			}
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}
