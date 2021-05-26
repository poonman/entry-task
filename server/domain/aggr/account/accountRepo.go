package account

type Repo interface {
	Get(username string) (a *Account, err error)
}
