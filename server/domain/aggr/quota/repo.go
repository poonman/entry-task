package quota

type Repo interface {
	Get(uid int64) (q *Quota, err error)
}
