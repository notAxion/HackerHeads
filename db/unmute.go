package db

import (
	"fmt"
	"strconv"
	"time"
)

type muteTime struct {
	GID        string    `db:"guild_id"`
	UserID     string    `db:"user_id"`
	UnmuteTime time.Time `db:"duration"`
}

// use natural join to get the role as well
func (db *DB) SaveUnmuteTime(gID, userID string, unmuteTime time.Time) error {
	_, err := db.Exec(`
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

func (db *DB) DeleteUnmuteTime(gID, userID int64) error {
	_, err := db.Exec(`
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

func (db *DB) GetMutedUsers() ([]muteTime, error) {
	users := make([]muteTime, 0, 10) //*todo inc cap if hosting
	err := db.Select(&users, `
		SELECT * FROM `+tableMuteTime+`
		;`)
	if err != nil {
		fmt.Println("Selcting all unmute error")
		return nil, err
	}
	return users, nil
}

type User struct {
	GID int64
	UID int64
}

func FromString(gID, uID string) User {
	g, _ := strconv.ParseInt(gID, 10, 64)
	u, _ := strconv.ParseInt(uID, 10, 64)
	usr := User{
		GID: g,
		UID: u,
	}
	return usr
}

func (usr *User) ToString() (gID, uID string) {
	gID = strconv.FormatInt(usr.GID, 10)
	uID = strconv.FormatInt(usr.UID, 10)
	return
}

type XMap map[User]time.Time

// take a XMap as a param
func (db *DB) TGetMutedUsers() (XMap, error) {
	users := make(XMap)
	rows, err := db.Queryx(`
	SELECT * FROM ` + tableMuteTime + `
	;`)
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for rows.Next() {
		var gID, uID int64
		var unmuteTime time.Time
		err := rows.Scan(&gID, &uID, &unmuteTime)
		if err != nil {
			fmt.Println("someone's scan had an error")
			return nil, err
		}
		k := User{
			gID,
			uID,
		}
		users[k] = unmuteTime
	}
	if err != nil {
		fmt.Println("unmute rows scan error")
		return nil, err
	}
	return users, err
}
