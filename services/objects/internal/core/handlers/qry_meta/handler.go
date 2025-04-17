package qry_meta

import (
	"context"
	"errors"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/common"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	repo   Repository
}

func NewMetaQryHandler(l *zap.Logger, repo Repository) *Handler {
	return &Handler{
		logger: l,
		repo:   repo,
	}
}

func (h *Handler) Handle(ctx context.Context, qry Query) (Response, error) {
	obj, err := qry.toObject()
	if err != nil {
		return Response{}, ErrInvalidRequestObjectID
	}

	err = h.repo.GetObjectByOID(ctx, obj)
	if err != nil {
		if errors.Is(err, common.ErrEmptyQueryResult) {
			return Response{}, nil
		}

		return Response{}, err
	}

	resp := newFromObject(*obj)
	return resp, nil
}
