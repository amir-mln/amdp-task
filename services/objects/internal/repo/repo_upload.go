package repo

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
)

var (
	_ cmd_upload.Repository = &DbRepository{}
)

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
	) VALUES (default, $1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT ON CONSTRAINT uq_objects_uid_name_mime
	DO UPDATE  SET "user_id" = EXCLUDED."user_id"
	RETURNING "id", "oid", "state";
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	tmp := entities.NewObject(0, "", "", nil)
	err = stmt.QueryRowContext(ctx, o.OID, o.UserID, o.Name, o.Mime, 0, "", o.State, o.CreatedAt).
		Scan(&tmp.ID, &tmp.OID, &tmp.State)
	if err == nil && tmp.OID.String() == o.OID.String() {
		o.ID = tmp.ID
	}
	if err == nil && tmp.OID.String() != o.OID.String() {
		*o = *tmp
		err = cmd_upload.ErrObjectExists
	}

	return err
}

func (r *DbRepository) SaveFinalObjTx(ctx context.Context, tx *sql.Tx, o *entities.Object) error {
	const query = `
		UPDATE "app"."objects"
		SET "size" = $1, "hash" = $2, "state" = $3
		WHERE "id" = $4;
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, o.Size, o.Hash, o.State, o.ID)
	return err
}
