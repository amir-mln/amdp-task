package messaging

import (
	"context"
	"database/sql"
	"time"

	"github.com/amir-mln/amdp-task/system/core/uow"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	uow.TxBeginner
	InsertMessage(context.Context, Message) error
	InsertMessageTx(context.Context, *sql.Tx, Message) error
	GetMessagesTx(context.Context, *sql.Tx, uint) ([]Message, error)
	DeleteMessageByIDTx(context.Context, *sql.Tx, uuid.UUID) error
}

type repo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRepository(logger *zap.Logger, db *sql.DB) Repository {
	return &repo{db: db, logger: logger}
}

func (r *repo) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *repo) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

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
		"user_id",
		"entity",
		"entity_id",
		"title",
		"type",
		"payload",
		"created_at",
		"publish_at"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,
		msg.ID, msg.UserId, msg.Entity, msg.EntityID,
		msg.Title, msg.Type, msg.Body, msg.CreatedAt, msg.publishAt,
	)
	return err
}

func (r *repo) GetMessagesTx(ctx context.Context, tx *sql.Tx, lim uint) ([]Message, error) {
	const query = `
	SELECT 
		"id",
		"user_id",
		"entity",
		"entity_id",
		"title",
		"type",
		"payload",
		"created_at"
	FROM "app"."messages"
	WHERE "publish_at" IS NULL OR "publish_at" <= NOW()
	LIMIT $1
	FOR UPDATE SKIP LOCKED
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx, lim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var m Message
		var userID, entityID sql.NullInt64
		var entity sql.NullString

		err := rows.Scan(
			&m.ID,
			&userID,
			&entity,
			&entityID,
			&m.Title,
			&m.Type,
			&m.Body,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			m.UserId = &userID.Int64
		}
		if entityID.Valid {
			m.EntityID = &entityID.Int64
		}
		if entity.Valid {
			m.Entity = &entity.String
		}

		m.PublishedAt = time.Now()
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *repo) DeleteMessageByIDTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	const query = `DELETE FROM "app"."messages" WHERE "id" = $1;`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
