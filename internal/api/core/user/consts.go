package user

import "time"

const (
	poolSize       = 10
	QuerySortField = "SortField"
	QueryLimit     = "Limit"
	QueryOffset    = "Offset"
	dataBasePort   = ":9091"
	topic          = "test_2808"
	eventCreate    = "Create"
	eventUpdate    = "Update"
	eventDelete    = "Delete"
	redisAddress   = "localhost:6379"
	redisDB        = 0
	redisPassword  = ""
	waitingtime    = time.Second * 10
	emptyID        = ""
)
