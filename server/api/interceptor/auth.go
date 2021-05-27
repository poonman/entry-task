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

var (
	ErrNoMetadataFound = status.New(status.BadRequest, "no metadata found")
	ErrNoUsernameFound = status.New(status.BadRequest, "no username found")
	ErrNoTokenFound = status.New(status.BadRequest, "no token found")
)

func (i *Interceptor) Do(ctx context.Context, in, out interface{}, serverInfo *server.InterceptorServerInfo, handler server.Handler) (err error) {

	err = i.Limit(ctx, serverInfo)
	if err != nil {
		return err
	}

	ctx, err = i.Auth(ctx, serverInfo)
	if err != nil {
		return err
	}

	return handler(ctx, in, out)
}

func (i *Interceptor) Auth(ctx context.Context, serverInfo *server.InterceptorServerInfo) (_ context.Context, err error) {
	if serverInfo.Method == "ReadSecureMessage" || serverInfo.Method == "WriteSecureMessage" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err = ErrNoMetadataFound
			return
		}

		username, ok := md["username"]
		if !ok {
			err = ErrNoUsernameFound
			return
		}
		token, ok := md["token"]
		if !ok {
			err = ErrNoTokenFound
			return
		}

		err = i.app.Authenticate(username, token)
		if err != nil {
			return
		}

		ctx = context.WithValue(ctx, "username", username)
		return ctx, nil
	}

	return ctx, nil
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
