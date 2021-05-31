package interceptor

import (
	"context"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/server/app"
	"strings"
)

type Interceptor struct {
	app *app.Service
}

var (
	ErrNoMetadataFound = status.New(status.BadRequest, "no metadata found")
	ErrNoUsernameFound = status.New(status.BadRequest, "no username found")
	ErrNoTokenFound    = status.New(status.BadRequest, "no token found")
	ErrNotAllowRead    = status.New(status.Unavailable, "not allow read due to rate limit")
	ErrNotAllowWrite   = status.New(status.Unavailable, "not allow write due to rate limit")
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

	if serverInfo.Method == "ReadSecureMessage" || serverInfo.Method == "WriteSecureMessage" {
		var (
			username string
			ok       bool
		)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err = ErrNoMetadataFound
			return
		}

		if username, ok = md["username"]; !ok {
			err = ErrNoUsernameFound
			return
		}

		if strings.Compare(serverInfo.Method, "ReadSecureMessage") == 0 {
			allow := i.app.AllowRead(ctx, username)
			if !allow {
				err = ErrNotAllowRead
				return
			}
		} else if strings.Compare(serverInfo.Method, "WriteSecureMessage") == 0 {
			allow := i.app.AllowWrite(ctx, username)
			if !allow {
				err = ErrNotAllowWrite
				return
			}
		}
	}
	return nil
}

func NewInterceptor(app *app.Service) (a *Interceptor) {
	return &Interceptor{
		app: app,
	}
}
