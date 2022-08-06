package config

import "time"

const (
	ApiKey          = "5471648580:AAFNRLmYMTC0HO--hwTZpdT8UL31zhAx09g"
	Host            = "localhost"
	Port            = 5432
	User            = "postgres"
	Password        = "admin"
	DBname          = "users"
	MaxConnIdleTime = time.Minute
	MaxConnLifetime = time.Hour
	MinConns        = 2
	MaxConns        = 4
)
