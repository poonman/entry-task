package account

import "context"

type Repo interface {
	Get(ctx context.Context, username string) (a *Account, err error)
}
