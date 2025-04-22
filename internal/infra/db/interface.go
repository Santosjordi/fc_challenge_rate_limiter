package db

type RateLimiterStorage interface {
	RegisterRequest(key string) error
	Increment(key string, windowSeconds int) (count int, err error)
	IsLockedOut(lockoutKey string) (bool, error)
	SetLockout(lockoutKey string, durationSeconds int) error
}
