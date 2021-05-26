package api

import (
	"context"
	"github.com/poonman/entry-task/server/idl/kv"
)

type Handler struct {
	//app *app.Service
	kv.UnimplementedStoreServer
}

func (h *Handler) WriteSecureMessage(ctx context.Context, req *kv.WriteSecureMessageReq, rsp *kv.WriteSecureMessageRsp) (err error) {
	//uid := uint64(0)
	//key := "key"
	//value := "value"
	//err = h.app.WriteSecureMessage(uid, key, value)
	//if err != nil {
	//	return
	//}

	return
}

func (h *Handler) ReadSecureMessage(ctx context.Context, req *kv.ReadSecureMessageReq, rsp *kv.ReadSecureMessageRsp) (err error) {
	//var (
	//	value string
	//)
	//uid := uint64(0)
	//
	//value, err = h.app.ReadSecureMessage(uid, req.Key)
	//if err != nil {
	//	return
	//}
	//
	//rsp.Value = value

	return
}

func NewHandler(
	//app *app.Service,
	) kv.StoreServer {

	h := &Handler{
		//app: app,
	}

	return h
}