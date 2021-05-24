package db

import (
	"database/sql"
	"fmt"
)

func (db *DB) MuteRoleID(gID string) (string, error) {
	var roleID string

	err := db.QueryRowx(`
		SELECT role_id from `+tableMute+` 
		where guild_id =$1
		;`, gID).Scan(&roleID)

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("scan error")
		}
		return "", err
	}

	return roleID, nil
}

func (db *DB) UpsertRole(gID, roleID string) error {
	_, err := db.Exec(`
		INSERT INTO `+tableMute+` (guild_id, role_id)
		VALUES
			($1, $2)
		ON CONFLICT (guild_id)
		DO
			UPDATE SET role_id = EXCLUDED.role_id
		;`, gID, roleID)

	if err != nil {
		return err
	}

	return nil
}
