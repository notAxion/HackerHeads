package db

import (
	"database/sql"
	"fmt"
)

func MuteRoleID(gID string) (string, error) {
	return PQ.MuteRoleID(gID)
}

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

func UpsertRole(gID, roleID string) error {
	return PQ.UpsertRole(gID, roleID)
}

func (db *DB) UpsertRole(gID, roleID string) error {
	res, err := db.Exec(`
		INSERT INTO `+tableMute+` (guild_id, role_id)
		VALUES
			($1, $2)
		ON CONFLICT (giuld_id)
		DO
			UPDATE SET role_id = EXCLUDED.role_id
		;`, gID, roleID)

	if err != nil {
		return err
	}
	fmt.Printf("%#v \n", res)

	return nil
}
