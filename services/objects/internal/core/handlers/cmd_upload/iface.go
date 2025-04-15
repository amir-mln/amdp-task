package cmd_upload

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"github.com/amir-mln/amdp-task/system/core/outbox"
)

type Repository interface {
	common.TxBeginner
	outbox.Repository
	GetExistingObjTx(context.Context, *sql.Tx, *entities.Object) (*entities.Object, error)
	SaveInitObjTx(context.Context, *sql.Tx, *entities.Object) error
	SaveFinalObj(context.Context, *entities.Object) error
}

type FileStore interface {
	PutObject(context.Context, *entities.Object) error
}
