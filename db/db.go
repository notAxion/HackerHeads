package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/notAxion/HackerHeads/constants"
	"github.com/notAxion/HackerHeads/secrets"
)

func init() {
	secrets.Set()

	dbName = os.Getenv(constants.DBName)
	tableMute = os.Getenv(constants.TableMute)
	tableMuteTime = os.Getenv(constants.TableMuteTime)
	tablePrefix = os.Getenv(constants.TablePrefix)
	dbUser = os.Getenv(constants.Username)
	dbPass = os.Getenv(constants.Pass)
	driverName = os.Getenv(constants.DriverName)

}

var (
	dbName, tableMute, tableMuteTime, tablePrefix, dbUser, dbPass, driverName string
)

type DB struct {
	*sqlx.DB
}

// NewDB will connect to Postgres // *todo add a recover maybe && change sslmode
func NewDB() *DB {
	defer fmt.Println("Database Set.")
	sourceName := fmt.Sprintf(`
		user=%s dbname=%s password=%s sslmode=disable
		`, dbUser, dbName, dbPass)
	return &DB{sqlx.MustConnect(driverName, sourceName)}
}
