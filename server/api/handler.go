package api

import (
	"context"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/app"
	"github.com/poonman/entry-task/server/domain/excp"
	"github.com/poonman/entry-task/server/idl/kv"
)

type Handler struct {
	app *app.Service
	kv.UnimplementedStoreServer
}

func (h *Handler) Login(ctx context.Context, req *kv.LoginReq, rsp *kv.LoginRsp) (err error) {

	rsp.Token, err = h.app.Login(ctx, req.Username, req.Password)
	if err != nil {
		if excp.IsBuzError(err) {
			rsp.Status = excp.Error2Status(err)
			err = nil
		}
		return
	}

	rsp.Status = excp.Error2Status(nil)

	return
}

func (h *Handler) WriteSecureMessage(ctx context.Context, req *kv.WriteSecureMessageReq, rsp *kv.WriteSecureMessageRsp) (err error) {

	username, ok := ctx.Value("username").(string)
	if !ok {
		err = status.ErrBadRequest
		return
	}

	err = h.app.WriteSecureMessage(username, req.Key, req.Value)
	if err != nil {
		return
	}

	rsp.Key = req.Key
	rsp.Value = req.Value

	return
}

func (h *Handler) ReadSecureMessage(ctx context.Context, req *kv.ReadSecureMessageReq, rsp *kv.ReadSecureMessageRsp) (err error) {

	username, ok := ctx.Value("username").(string)
	if !ok {
		err = status.ErrBadRequest
		return
	}

	var (
		value string
	)

	value, err = h.app.ReadSecureMessage(username, req.Key)
	if err != nil {
		if excp.IsBuzError(err) {
			rsp.Status = excp.Error2Status(err)
			err = nil
		}
		return
	}

	rsp.Value = value

	return
}

func NewHandler(
	app *app.Service,
) kv.StoreServer {

	h := &Handler{
		app: app,
	}

	return h
}
