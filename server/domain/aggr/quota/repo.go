package quota

import "context"

type Repo interface {
	Get(ctx context.Context, username string) (q *Quota, err error)
}
