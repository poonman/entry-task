package interceptor

import (
	"context"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/app"
)

type AuthInterceptor struct {
	app *app.Service
}

func (a *AuthInterceptor) Auth(ctx context.Context, in, out interface{}, serverInfo *server.InterceptorServerInfo, handler server.Handler) (err error) {
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

		err = a.app.Authenticate(username, token)
		if err != nil {
			return
		}

		return
	}

	return handler(ctx, in, out)
}

func NewAuthInterceptor(app *app.Service) (a *AuthInterceptor) {
	return &AuthInterceptor{
		app: app,
	}
}
