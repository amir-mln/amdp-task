package common

import (
	"context"
	"database/sql"
)

type TxBeginner interface {
	Begin() (*sql.Tx, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}
