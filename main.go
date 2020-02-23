package main

import (
	"context"
	"encoding/json"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/trewanek/transaction-common-func/db/rdb"
	"github.com/trewanek/transaction-common-func/presenter"
)

func main() {
	ctx := context.Background()
	dbConn, err := rdb.NewDBConn()
	if err != nil {
		log.Fatal(err)
	}

	q1 := `INSERT INTO users(user_name, email, telephone) VALUES('hoge', 'hoge@hoge.com', 'XXX-XXXX-XXXX');`
	q2 := `SELECT * FROM users;`

	err = rdb.Transact(ctx, dbConn, func(tx *sqlx.Tx) (err error) {
		if _, err = tx.Exec(q1); err != nil {
			return err
		}
		rows, err := tx.Queryx(q2)
		defer func() {
			rows.Close()
		}()
		if err != nil {
			return err
		}

		list := make([]map[string]interface{}, 0, 0)
		for rows.Next() {
			dst := map[string]interface{}{}
			err = rows.MapScan(dst)
			if err != nil {
				break
			}
			list = append(list, dst)
		}

		bs, err := json.Marshal(list)
		if err != nil {
			return err
		}

		pre := presenter.NewStdoutPresenter()
		_, err = pre.Write(bs)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
