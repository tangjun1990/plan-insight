package kredis

import "github.com/go-redis/redis/v8"

type Err string

func (e Err) Error() string { return string(e) }

const (
	ErrInvalidParams = Err("invalid params")

	ErrNotObtained = Err("redislock: not obtained")

	ErrLockNotHeld = Err("redislock: lock not held")

	//Nil reply returned by Redis when key does not exist.
	Nil = redis.Nil
)
