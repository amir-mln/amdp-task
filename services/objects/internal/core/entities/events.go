package entities

type ObjectUploadCompleted struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	ObjID  string `json:"object_id"`
	Name   string `json:"name"`
	Mime   string `json:"mime"`
	Size   uint64 `json:"size"`
	Hash   string `json:"hash"`
	State  string `json:"state"`
}

type ObjectUploadFailed struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	ObjID  string `json:"object_id"`
	Name   string `json:"name"`
	Mime   string `json:"mime"`
}

type IncompleteObjectInserted struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	ObjID  string `json:"object_id"`
	Name   string `json:"name"`
	Mime   string `json:"mime"`
	State  string `json:"state"`
}
