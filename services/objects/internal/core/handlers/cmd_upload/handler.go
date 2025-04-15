package cmd_upload

import (
	"context"
	"errors"
	"time"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"github.com/amir-mln/amdp-task/system/core/outbox"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	logger    *zap.Logger
	repo      Repository
	filestore FileStore
}

func NewUploadCmdHandler(l *zap.Logger, repo Repository, fs FileStore) *Handler {
	return &Handler{
		logger:    l,
		repo:      repo,
		filestore: fs,
	}
}

func (h *Handler) saveIncompleteObj(ctx context.Context, obj *entities.Object) (err error) {
	tx, err := h.repo.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = h.repo.SaveInitObjTx(ctx, tx, obj)
	if errors.Is(err, ErrObjectExists) {
		var tmp *entities.Object
		tmp, err = h.repo.GetExistingObjTx(ctx, tx, obj)
		if errors.Is(err, common.ErrEmptyQueryResult) {
			//TODO: wrap err in custom err of type internal
		} else if err == nil {
			*obj = *tmp
		}
	}
	if err != nil {
		return
	}

	rec := outbox.Record{
		UUID:     uuid.New(),
		Entity:   "Object",
		EntityID: obj.ID,
		Title:    "",
		Target:   "object.outgoing",
		Payload: entities.IncompleteObjectInserted{
			ID:     obj.ID,
			UserID: obj.UserID,
			ObjID:  obj.OID.String(),
			Name:   obj.Name,
			Mime:   obj.Mime,
		},
		CreatedAt: time.Now(),
	}
	err = h.repo.InsertRecordTx(ctx, tx, rec)
	return
}

func (h *Handler) uploadObject(ctx context.Context, obj *entities.Object) (err error) {
	rec := outbox.Record{
		UUID:      uuid.New(),
		Entity:    "Object",
		EntityID:  obj.ID,
		Title:     "",
		Target:    "object.outgoing",
		Payload:   nil,
		CreatedAt: time.Now(),
	}
	if err = h.filestore.PutObject(ctx, obj); err != nil {
		rec.Title = "ObjectUploadFailed"
		rec.Payload = entities.ObjectUploadFailed{
			ID:     obj.ID,
			ObjID:  obj.OID.String(),
			UserID: obj.UserID,
			Name:   obj.Name,
			Mime:   obj.Mime,
		}
		obj.State = entities.Failed
	} else {
		rec.Title = "ObjectUploadCompleted"
		rec.Payload = entities.ObjectUploadCompleted{
			UserID: obj.UserID,
			ObjID:  obj.OID.String(),
			Name:   obj.Name,
			Mime:   obj.Mime,
			Size:   obj.Size,
			Hash:   obj.Hash,
		}
		obj.State = entities.Complete
	}

	sErr := h.repo.SaveFinalObj(ctx, obj)
	if sErr != nil {
		// log error
	}
	rErr := h.repo.InsertRecord(ctx, rec)
	if rErr != nil {
		// log error
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (Response, error) {
	obj := entities.NewObject(cmd.UserID, cmd.Mime, cmd.Name, cmd.Object)
	err := h.saveIncompleteObj(ctx, obj)
	if err != nil {
		if errors.Is(err, ErrObjectExists) {
			return Response{OID: obj.OID.String()}, nil
		}

		return Response{}, err
	}

	err = h.uploadObject(ctx, obj)
	h.logger.Info("Final object", zap.String("Hash", obj.Hash), zap.Uint64("Size", obj.Size))
	return Response{OID: obj.OID.String()}, nil
}
