package secrets

import "os"

func Set() {
	os.Setenv("TABLE_MUTE", "")
	os.Setenv("TABLE_MUTE_TIME", "")
	os.Setenv("PQ_DB_NAME", "")
	os.Setenv("PQ_USERNAME", "")
	os.Setenv("PQ_PASSWORD", "")
	os.Setenv("SQL_DRIVER", "")
}
