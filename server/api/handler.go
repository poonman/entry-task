package api

import (
	"context"
	"github.com/poonman/entry-task/server/app"
)

type Handler struct {
	app *app.Service
}

func (h *Handler) WriteSecureMessage(ctx context.Context, req interface{}, rsp interface{}) (err error) {
	uid := uint64(0)
	key := "key"
	value := "value"
	err = h.app.WriteSecureMessage(uid, key, value)
	if err != nil {
		return
	}

	return
}

func (h *Handler) ReadSecureMessage(ctx context.Context, req interface{}, rsp interface{}) (err error) {
	var (
		value string
	)
	uid := uint64(0)
	key := "key"
	value, err = h.app.ReadSecureMessage(uid, key)
	if err != nil {
		return
	}

	rsp = value

	return
}
