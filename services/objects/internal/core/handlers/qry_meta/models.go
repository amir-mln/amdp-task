package qry_meta

import (
	"github.com/amir-mln/amdp-task/services/objects/internal/core/entities"
	"github.com/google/uuid"
)

type Query struct {
	// Should be parsed from JWT and read from HTTP Request's context
	// It's not supported in the current version
	UserID int64
	OID    string
}

func (req Query) toObject() (*entities.Object, error) {
	oid, err := uuid.Parse(req.OID)
	if err != nil {
		return nil, ErrInvalidRequestObjectID.WithArgs(req.OID)
	}

	obj := entities.NewObject(req.UserID, "", "", nil)
	obj.OID = oid
	return obj, nil
}

type Response struct {
	ObjID string `json:"object_id"`
	Name  string `json:"name"`
	Mime  string `json:"mime"`
	Size  int64  `json:"size"`
	Hash  string `json:"hash"`
	State string `json:"state"`
}

func newFromObject(obj entities.Object) Response {
	resp := Response{
		ObjID: obj.OID.String(),
		Name:  obj.Name,
		Mime:  obj.Mime,
		Size:  obj.Size,
		Hash:  obj.Hash,
		State: obj.State.String(),
	}

	return resp
}
