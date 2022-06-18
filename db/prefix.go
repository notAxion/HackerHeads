package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *DB) GetAllPrefix(prefixMap *map[string]string) (string, error) {

	rows, err := db.Queryx(`
		SELECT * from ` + tablePrefix + `
		;`)

	for rows.Next() {

	}
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("scan error")
		}
		return "", err
	}

	return "", errors.New("implement GetAllPrefix") // return prefix, nil
}

func (db *DB) UpsertPrefix(gID, prefix string) error {
	_, err := db.Exec(`
		INSERT INTO `+tablePrefix+` (guild_id, prefix)
		VALUES
			($1, $2)
		ON CONFLICT (guild_id)
		DO
			UPDATE SET prefix = EXCLUDED.prefix
		;`, gID, prefix)

	if err != nil {
		fmt.Println("upsert prefix err")
		return err
	}

	return nil
}
