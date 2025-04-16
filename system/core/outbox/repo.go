package outbox

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type Repository interface {
	InsertRecord(context.Context, Record) error
	InsertRecordTx(context.Context, *sql.Tx, Record) error
}

type repo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger, db *sql.DB) Repository {
	return &repo{db: db, logger: logger}
}

// InsertRecordTx implements Repository.
func (r *repo) InsertRecord(ctx context.Context, rec Record) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = r.InsertRecordTx(ctx, tx, rec)
	return
}

func (r *repo) InsertRecordTx(ctx context.Context, tx *sql.Tx, rec Record) error {
	const query = `
	INSERT INTO "app"."outbox" (
		"id",
		"uuid",
		"entity",
		"entity_id",
		"title",
		"target",
		"payload",
		"created_at",
		"processed_at"
	) VALUES (default, $1, $2, $3, $4, $5, $6, $7, NULL);
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rec.UUID, rec.Entity, rec.EntityID, rec.Title, rec.Target, rec.Payload, rec.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
