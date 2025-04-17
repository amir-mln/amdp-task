package messaging

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type Repository interface {
	InsertMessage(context.Context, Message) error
	InsertMessageTx(context.Context, *sql.Tx, Message) error
}

type repo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger, db *sql.DB) Repository {
	return &repo{db: db, logger: logger}
}

// InsertRecordTx implements Repository.
func (r *repo) InsertMessage(ctx context.Context, ob Message) (err error) {
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
	err = r.InsertMessageTx(ctx, tx, ob)
	return
}

func (r *repo) InsertMessageTx(ctx context.Context, tx *sql.Tx, msg Message) error {
	const query = `
	INSERT INTO "app"."messages" (
		"id",
		"tx_id",
		"user_id",
		"entity",
		"entity_id",
		"title",
		"type",
		"payload",
		"created_at",
		"publish_at"
	) VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,
		msg.Header.TxID, msg.Header.UserId, msg.Header.Entity, msg.Header.EntityID,
		msg.Header.Title, msg.Header.Type, msg.Body, msg.Header.CreatedAt, msg.publishAt,
	)
	return err
}
