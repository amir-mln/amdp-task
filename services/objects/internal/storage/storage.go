package storage

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

const bucketName = "def"

type ObjStorage struct {
	logger *zap.Logger
	client *minio.Client
}

func NewObjectStorage(logger *zap.Logger, c *minio.Client) *ObjStorage {
	return &ObjStorage{
		logger: logger,
		client: c,
	}
}
