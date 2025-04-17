package cmd_upload

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/amir-mln/amdp-task/system/core/messaging"
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

func (h *Handler) saveInitialObj(ctx context.Context, obj *entities.Object, msg messaging.Message) (err error) {
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
	if err != nil {
		return
	}

	err = h.repo.InsertMessageTx(ctx, tx, msg)
	return
}

func (h *Handler) saveFinalObj(ctx context.Context, obj *entities.Object, msg messaging.Message) {
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

	err = h.repo.SaveFinalObjTx(ctx, tx, obj)
	if err != nil {
		opts := []zap.Field{
			zap.Stringer("Object UUID", obj.OID),
			zap.Stringer("State", obj.State),
			zap.NamedError("Cause", err),
		}
		h.logger.Warn("Failed to save the final state of object upload", opts...)
		return
	}
	err = h.repo.InsertMessageTx(ctx, tx, msg)
	if err != nil {
		b, _ := json.Marshal(msg.Body)
		opts := []zap.Field{
			zap.Stringer("Object UUID", obj.OID),
			zap.ByteString("Payload", b),
			zap.NamedError("Cause", err),
		}
		h.logger.Warn("Failed to save the message of final object upload", opts...)
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (Response, error) {
	obj := entities.NewObject(cmd.UserID, cmd.Mime, cmd.Name, cmd.Object)

	eve := &entities.InitialObjectInserted{
		ID:     obj.ID,
		UserID: obj.UserID,
		ObjID:  obj.OID,
	}
	msg := messaging.NewMessage(eve, messaging.WithEntity(obj))
	err := h.saveInitialObj(ctx, obj, *msg)
	if errors.Is(err, ErrObjectExists) {
		return Response{OID: obj.OID, State: obj.State.String()}, nil
	}
	if err != nil {
		return Response{}, err
	}

	err = h.filestore.PutObject(ctx, obj)
	fmsg, resp := &messaging.Message{}, Response{}
	if err != nil {
		obj.State = entities.Failed
		eve := &entities.ObjectUploadFailed{
			ID:     obj.ID,
			ObjID:  obj.OID,
			UserID: obj.UserID,
			Error:  err.Error(),
		}
		fmsg = messaging.NewMessage(eve, messaging.WithEntity(obj))
	} else {
		obj.State = entities.Completed
		eve := &entities.ObjectUploadCompleted{
			UserID: obj.UserID,
			ObjID:  obj.OID,
			Name:   obj.Name,
			Mime:   obj.Mime,
			Size:   obj.Size,
			Hash:   obj.Hash,
		}
		fmsg = messaging.NewMessage(eve, messaging.WithEntity(obj))
		resp.OID = obj.OID
		resp.State = obj.State.String()
	}

	go h.saveFinalObj(context.Background(), obj, *fmsg)
	return resp, err
}
