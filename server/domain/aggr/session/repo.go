package session

type Repo interface {
	Save(username, token string) (err error)
	Get(username string) (token string, err error)
}
