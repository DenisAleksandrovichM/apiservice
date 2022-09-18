package config

const (
	DBHost       = "localhost"
	DBPort       = 5432
	DBUser       = "postgres"
	DBPassword   = "admin"
	DBName       = "users"
	GRPCNetwork  = "tcp"
	GRPCAddress  = ":9091"
	HTTPEndpoint = ":9091"
	HTTPAddress  = ":9090"
)

var (
	Brokers = []string{"localhost:19091"}
	Topics  = []string{"test_2808"}
)
