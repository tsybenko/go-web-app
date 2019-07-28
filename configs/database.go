package configs

import "os"

// Gets "env" value from environment variables if exists.
// Otherwise returns "val" as default value
func envOrVal(env string, val string) string {
	if v := os.Getenv(env); v != "" {
		return v
	} else {
		return val
	}
}

var DbConfig = map[string]func() string{
	"username": func() string {
		return envOrVal("DB_USERNAME", "root")
	},
	"password": func() string {
		return envOrVal("DB_PASS", "")
	},
	"host": func() string {
		return envOrVal("DB_HOST", "127.0.0.1")
	},
	"port": func() string {
		return envOrVal("DB_PORT", "3306")
	},
	"dbname": func() string {
		return envOrVal("DB_NAME", "")
	},
}