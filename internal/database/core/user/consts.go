package user

import "time"

const (
	poolSize        = 10
	tableName       = "users"
	tableColumns    = "login, first_name, last_name, weight, height, age"
	primaryKey      = "login"
	QuerySortField  = "SortField"
	QueryLimit      = "Limit"
	QueryOffset     = "Offset"
	redisAddress    = "localhost:6379"
	redisDB         = 0
	redisPassword   = ""
	redisExpiration = 2 * time.Hour
)
