package user

import "github.com/pkg/errors"

var (
	errUserNotExists = errors.New("user does not exist")
	errUserExists    = errors.New("user exists")
	errManyUsers     = errors.New("there are many users with this login")
	errSQL           = errors.New("SQL error")
	errAdd           = errors.New("add error")
	errRead          = errors.New("read error")
	errUpdate        = errors.New("update error")
	errDelete        = errors.New("delete error")
	errList          = errors.New("list error")
	errRedisGet      = errors.New("redis on get error")
	errRedisScan     = errors.New("redis on scan error")
)
