package repo

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/system/core/uow"
)

var _ uow.TxBeginner = &DbRepository{}

func (r *DbRepository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *DbRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}
