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
	dbUser = os.Getenv(constants.Username)
	dbPass = os.Getenv(constants.Pass)
	driverName = os.Getenv(constants.DriverName)

	PQ = NewDB() // prepareing bindvars is $n
}

var (
	dbName, tableMute, tableMuteTime, dbUser, dbPass, driverName string
	PQ                                                           *DB
)

type DB struct {
	*sqlx.DB
}

// NewDB will connect to Postgres // *todo add a recover maybe
func NewDB() *DB {
	defer fmt.Println("Database Set.")
	sourceName := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " sslmode=disable"
	return &DB{sqlx.MustConnect(driverName, sourceName)}
}
