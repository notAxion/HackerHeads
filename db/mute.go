package db

import (
	"database/sql"
	"fmt"
)

// try both parsing to int and not to see what works
func (db *DB) MuteRoleID(idStr string) (string, error) {
	var roleID string

	err := db.QueryRowx(`
		SELECT role_id from `+tableMute+` 
		where guild_id =$1
		;`, idStr).Scan(&roleID)

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("scan error")
		}
		return "", err
	}

	return roleID, nil
}
