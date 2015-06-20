package redisq

import (
	"errors"
	"strconv"

	"github.com/fortytw2/embercrest/user"
	"github.com/garyburd/redigo/redis"
)

// Queue is a concurrency safe queue
type Queue struct {
	r *redis.Pool
}

func NewQueue(r *redis.Pool) *Queue {
	return &Queue{
		r: r,
	}
}

func (rq *Queue) Enter(u *user.User) error {
	numAdded, err := redis.Int64(rq.r.Get().Do("ZADD", "queue", u.Elo, u.Username))
	if err != nil {
		return err
	}
	if numAdded == 0 {
		return errors.New("user already in queue")
	}

	if numAdded != 1 {
		return errors.New("user not added to queue ")
	}

	return nil
}

func (rq *Queue) Exit(username string) error {
	numRem, err := redis.Int64(rq.r.Get().Do("ZREM", "queue", username))
	if err != nil {
		return err
	}
	if numRem != 1 {
		return errors.New("user not in queue, so cannot be removed")
	}
	return nil
}

func (rq *Queue) Pop() (username string, err error) {
	var users []string
	users, err = redis.Strings(rq.r.Get().Do("ZRANGE", "queue", "0", "0"))
	if err != nil {
		return
	}

	if len(users) != 1 {
		err = errors.New("something went wrong with redis")
		return
	}

	username = users[0]

	rq.r.Get().Do("ZREM", "queue", users[0])

	return
}

func (rq *Queue) Length() (length int64, err error) {
	length, err = redis.Int64(rq.r.Get().Do("ZCOUNT", "queue", "-inf", "+inf"))
	if err != nil {
		return
	}
	return
}

func (rq *Queue) String() string {
	length, _ := rq.Length()

	if length >= 2 {
		return strconv.FormatInt(length, 10) + " players in queue"
	} else if length == 1 {
		return strconv.FormatInt(length, 10) + " player in queue"
	} else {
		return "no players in queue"
	}
}
