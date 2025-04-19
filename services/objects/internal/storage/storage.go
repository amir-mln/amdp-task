package storage

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type ObjStorage struct {
	logger *zap.Logger
	client *minio.Client
	Bucket string
}

func NewObjectStorage(logger *zap.Logger, c *minio.Client, buc string) *ObjStorage {
	return &ObjStorage{
		logger: logger,
		client: c,
		Bucket: buc,
	}
}
