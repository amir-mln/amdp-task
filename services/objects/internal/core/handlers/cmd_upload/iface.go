package cmd_upload

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/system/core/messaging"
	"github.com/amir-mln/amdp-task/system/core/uow"
)

type Repository interface {
	uow.TxBeginner
	messaging.Repository
	// GetExistingObj(context.Context, *entities.Object) (*entities.Object, error)
	SaveInitObjTx(context.Context, *sql.Tx, *entities.Object) error
	SaveFinalObjTx(context.Context, *sql.Tx, *entities.Object) error
}

type FileStore interface {
	PutObject(context.Context, *entities.Object) error
}
