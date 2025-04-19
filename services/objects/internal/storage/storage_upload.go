package storage

import (
	"context"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
	"github.com/minio/minio-go/v7"
)

var _ cmd_upload.FileStore = &ObjStorage{}

func (os *ObjStorage) PutObject(ctx context.Context, o *entities.Object) error {
	_, err := os.client.PutObject(
		ctx,
		os.Bucket,
		o.OID.String(),
		o,
		-1,
		minio.PutObjectOptions{ContentType: o.Mime},
	)
	return err
}
