package ratelimiter

type Repo interface {
	Get(username string) (lim *RateLimiter)
	Save(lim *RateLimiter)
}
