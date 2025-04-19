package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/qry_meta"
)

var _ qry_meta.Repository = &DbRepository{}

func (repo *DbRepository) GetObjectByOID(ctx context.Context, obj *entities.Object) error {
	const query = ` 
	SELECT "oid", "name", "mime", "size", "hash", "state"
	FROM "app"."objects"
	WHERE "user_id" = $1 AND "oid" = $2;
	`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	row := stmt.QueryRowContext(ctx, obj.UserID, obj.OID)
	if err := row.Scan(&obj.OID, &obj.Name, &obj.Mime, &obj.Size, &obj.Hash, &obj.State); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.ErrEmptyQueryResult
		}

		return err
	}

	return nil
}
