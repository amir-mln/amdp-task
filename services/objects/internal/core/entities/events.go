package entities

import (
	"github.com/amir-mln/amdp-task/system/core/messaging"
	"github.com/google/uuid"
)

type InitialObjectInserted struct {
	ID     int64     `json:"id"`
	UserID int64     `json:"user_id"`
	ObjID  uuid.UUID `json:"object_id"`
	State  string    `json:"state"`
}

func (i InitialObjectInserted) MessageTitle() string {
	return "InitialObjectInserted"
}

func (i InitialObjectInserted) MessageType() messaging.MessageType {
	return messaging.Event
}

type ObjectUploadCompleted struct {
	ID     int64     `json:"id"`
	UserID int64     `json:"user_id"`
	ObjID  uuid.UUID `json:"object_id"`
	Name   string    `json:"name"`
	Mime   string    `json:"mime"`
	Size   int64     `json:"size"`
	Hash   string    `json:"hash"`
	State  string    `json:"state"`
}

func (o ObjectUploadCompleted) MessageTitle() string {
	return "ObjectUploadCompleted"
}

func (o ObjectUploadCompleted) MessageType() messaging.MessageType {
	return messaging.Event
}

type ObjectUploadFailed struct {
	ID     int64     `json:"id"`
	UserID int64     `json:"user_id"`
	ObjID  uuid.UUID `json:"object_id"`
	Name   string    `json:"name"`
	Mime   string    `json:"mime"`
	Error  string    `json:"error"`
	State  string    `json:"state"`
}

func (o ObjectUploadFailed) MessageTitle() string {
	return "ObjectUploadFailed"
}

func (o ObjectUploadFailed) MessageType() messaging.MessageType {
	return messaging.Event
}
