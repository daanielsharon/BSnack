package common

import (
	"bsnack/app/pkg/config"
	"fmt"
)

func GenerateDSN(cfg config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
		cfg.Postgres.Timezone,
	)
}
