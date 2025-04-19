package objects

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/cmd_upload"
	"github.com/amir-mln/amdp-task/services/objects/internal/core/handlers/qry_meta"
	"github.com/amir-mln/amdp-task/system/core/bus"
	syserr "github.com/amir-mln/amdp-task/system/errors"
	"go.uber.org/zap"
)

type ObjectsRouter struct {
	logger *zap.Logger
	bus    *bus.HandlerBus
}

func NewObjectRouter(l *zap.Logger, b *bus.HandlerBus) *ObjectsRouter {
	return &ObjectsRouter{
		logger: l,
		bus:    b,
	}
}

func (or *ObjectsRouter) Router() http.Handler {
	mux := &http.ServeMux{}
	mux.HandleFunc("PUT /objects/{$}", or.HandlePutObject)
	mux.HandleFunc("GET /objects/{objectid}/meta/{$}", or.HandleGetObjectMeta)
	return mux
}

func (router *ObjectsRouter) HandlePutObject(w http.ResponseWriter, r *http.Request) {
	const expectedFormName = "file"
	mr, err := r.MultipartReader()
	if err != nil {
		msg := fmt.Sprintf("Invalid multipart form; %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	part, err := mr.NextPart()
	if err != nil {
		msg := fmt.Sprintf("Error reading multipart data; %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if fn := part.FormName(); fn != expectedFormName {
		msg := fmt.Sprintf("invalid multipart form name. Expected %q, got %q", expectedFormName, fn)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	req := cmd_upload.Command{
		UserID: 0,
		Object: part,
		Name:   part.FileName(),
		Mime:   part.Header.Get("Content-Type"),
	}
	resp, err := bus.Handle[cmd_upload.Response](r.Context(), router.bus, req)
	if err != nil {
		syserr.HandleHTTPError(w, err, router.logger)
		return
	}

	b, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		syserr.HandleHTTPError(w, err, router.logger)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (router *ObjectsRouter) HandleGetObjectMeta(w http.ResponseWriter, r *http.Request) {
	req := qry_meta.Query{UserID: 0, OID: r.PathValue("objectid")}
	resp, err := bus.Handle[qry_meta.Response](r.Context(), router.bus, req)
	if err != nil {
		syserr.HandleHTTPError(w, err, router.logger)
		return
	}

	b, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		syserr.HandleHTTPError(w, err, router.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}
