package qry_meta

import (
	"context"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
)

type Repository interface {
	GetObjectByOID(context.Context, *entities.Object) error
}
