package db

import (
	"fmt"
	"time"
)

type muteTime struct {
	GID        string    `db:"guild_id"`
	UserID     string    `db:"user_id"`
	UnmuteTime time.Time `db:"duration"`
}

// use natural join to get the role as well
func SaveUnmuteTime(gID, userID string, unmuteTime time.Time) error {
	_, err := PQ.Exec(`
		INSERT INTO `+tableMuteTime+` (guild_id, user_id, duration)
		VALUES
		($1, $2, $3)
		;`, gID, userID, unmuteTime)
	if err != nil {
		fmt.Println("Insert unmute error")
		return err
	}
	return nil
}

func DeleteUnmuteTime(gID, userID string) error {
	_, err := PQ.Exec(`
		DELETE FROM `+tableMuteTime+` 
		WHERE guild_id=$1 
			AND user_id=$2
		;`, gID, userID)
	if err != nil {
		fmt.Println("Delete unmute error")
		return err
	}
	return nil
}

func GetMutedUsers() ([]muteTime, error) {
	users := make([]muteTime, 0, 10) //*todo inc cap if hosting
	err := PQ.Select(&users, `
		SELECT * FROM `+tableMuteTime+`
		;`)
	if err != nil {
		fmt.Println("Selcting all unmute error")
		return nil, err
	}
	return users, nil
}

type XMap map[string]time.Time

// take a XMap as a param
func TGetMutedUsers() (XMap, error) {
	users := make(XMap)
	rows, err := PQ.Queryx(`
	SELECT * FROM ` + tableMuteTime + `
		;`)
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for rows.Next() {
		usr := muteTime{}
		rows.StructScan(&usr)
		k := usr.GID + " " + usr.UserID
		users[k] = usr.UnmuteTime
	}
	if err != nil {
		fmt.Println("unmute rows scan error")
		return nil, err
	}
	return users, err
}
