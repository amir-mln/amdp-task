package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"github.com/jackc/pgconn"
)

var (
	_ cmd_upload.Repository = &DbRepository{}
)

const (
	objectUniqueConstraint = "uq_objects_uid_name_mime"
)

func (r *DbRepository) GetExistingObjTx(ctx context.Context, tx *sql.Tx, obj *entities.Object) (*entities.Object, error) {
	const query = `
		SELECT "oid" FROM "app"."objects"
		WHERE "user_id" = $1 AND "name" = $2 AND "mime" = $3;
	`
	var o entities.Object

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, query, obj.UserID, obj.Name, obj.Mime)
	if err := row.Scan(&o.ID, &o.OID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrEmptyQueryResult
		}
		return nil, err
	}

	return &o, nil
}

func (r *DbRepository) SaveInitObjTx(ctx context.Context, tx *sql.Tx, o *entities.Object) error {
	const query = `
	INSERT INTO "app"."objects" (
		"id",
		"oid",
		"user_id",
		"name",
		"mime",
		"size",
		"hash",
		"state",
		"created_at"
	) VALUES (default, $1, $2, $3, $4, $5, 0, $6, $7)
	RETURNING "id";
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	err = stmt.QueryRowContext(ctx, o.ID, o.OID, o.UserID, o.Name, o.Mime, o.State, o.CreatedAt).
		Scan(&o.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == objectUniqueConstraint {
			return cmd_upload.ErrObjectExists
		}
	}

	return nil
}

func (r *DbRepository) SaveFinalObj(ctx context.Context, o *entities.Object) error {
	const query = `
		UPDATE "app"."objects"
		SET "size" = $1, "hash" = $2, "state" = $3
		WHERE "id" = $4;
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, query, o.Size, o.Hash, o.State, o.ID)
	return err
}
