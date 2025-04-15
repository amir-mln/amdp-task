package repo

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
)

var _ common.TxBeginner = &DbRepository{}

func (r *DbRepository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *DbRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}
