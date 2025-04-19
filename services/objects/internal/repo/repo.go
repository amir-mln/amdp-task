package repo

import (
	"database/sql"

	"github.com/amir-mln/amdp-task/system/core/messaging"
	"go.uber.org/zap"
)

type DbRepository struct {
	logger *zap.Logger
	db     *sql.DB
	messaging.Repository
}

func NewDbRepository(logger *zap.Logger, db *sql.DB) *DbRepository {
	or := messaging.NewRepository(logger, db)
	return &DbRepository{
		logger:     logger,
		db:         db,
		Repository: or,
	}
}
