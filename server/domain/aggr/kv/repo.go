package kv

type Repo interface {
	Set(key, value string) (err error)
	Get(key string) (value string, err error)
}
