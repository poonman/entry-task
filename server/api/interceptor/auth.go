package interceptor

import (
	"context"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/app"
)

type Interceptor struct {
	app *app.Service
}

func (i *Interceptor) Do(ctx context.Context, in, out interface{}, serverInfo *server.InterceptorServerInfo, handler server.Handler) (err error) {

	err = i.Limit(ctx, serverInfo)
	if err != nil {
		return err
	}

	err = i.Auth(ctx, serverInfo)
	if err != nil {
		return err
	}

	return handler(ctx, in, out)
}

func (i *Interceptor) Auth(ctx context.Context, serverInfo *server.InterceptorServerInfo) (err error) {
	if serverInfo.Method == "ReadSecureMessage" || serverInfo.Method == "WriteSecureMessage" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err = status.ErrBadRequest
			return
		}

		username, ok := md["username"]
		if !ok {
			err = status.ErrBadRequest
			return
		}
		token, ok := md["token"]
		if !ok {
			err = status.ErrBadRequest
			return
		}

		err = i.app.Authenticate(username, token)
		if err != nil {
			return
		}

		return
	}

	return nil
}

func (i *Interceptor) Limit(ctx context.Context, serverInfo *server.InterceptorServerInfo) (err error) {
	// todo:
	return nil
}

func NewInterceptor(app *app.Service) (a *Interceptor) {
	return &Interceptor{
		app: app,
	}
}
