package cmd_upload

import (
	"context"
	"database/sql"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"github.com/amir-mln/amdp-task/system/core/messaging"
)

type Repository interface {
	common.TxBeginner
	messaging.Repository
	// GetExistingObj(context.Context, *entities.Object) (*entities.Object, error)
	SaveInitObjTx(context.Context, *sql.Tx, *entities.Object) error
	SaveFinalObjTx(context.Context, *sql.Tx, *entities.Object) error
}

type FileStore interface {
	PutObject(context.Context, *entities.Object) error
}
