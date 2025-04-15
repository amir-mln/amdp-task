package repo

import (
	"database/sql"

	"github.com/amir-mln/amdp-task/system/core/outbox"
	"go.uber.org/zap"
)

type DbRepository struct {
	logger *zap.Logger
	db     *sql.DB
	outbox.Repository
}

func NewDbRepository(logger *zap.Logger, db *sql.DB) *DbRepository {
	or := outbox.NewRepository(logger, db)
	return &DbRepository{
		logger:     logger,
		db:         db,
		Repository: or,
	}
}
