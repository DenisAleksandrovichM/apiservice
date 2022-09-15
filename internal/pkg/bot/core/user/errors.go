package user

import "github.com/pkg/errors"

var (
	errAdd           = errors.New("add error")
	errRead          = errors.New("read error")
	errUpdate        = errors.New("update error")
	errDelete        = errors.New("delete error")
	errList          = errors.New("list error")
	errDataBase      = errors.New("database error")
	errKafka         = errors.New("kafka error")
	errRedisGet      = errors.New("redis on get error")
	errRedisScan     = errors.New("redis on scan error")
	errCheckCUResult = errors.New("check create/update result error")
	errMarshal       = errors.New("on marshal error")
)
