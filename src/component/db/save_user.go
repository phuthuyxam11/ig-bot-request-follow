package db

import (
	"database/sql"
	"fmt"
)

type UserModel struct {
	PkId       string
	UserName   string
	IsPrivate  int
	FollowFlag string
	CreatedAt  string
	ModifiedAt string
}

func SaveUser(db *sql.DB, userModels []UserModel) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// execute
	insertValue := ""
	for index, user := range userModels {
		separateCharset := ""
		if index != len(userModels)-1 {
			separateCharset = ","
		}
		isExists, _ := CheckUserExists(db, user.PkId)
		if !isExists {
			insertValue += fmt.Sprintf(`('%s','%s',%d,'%s','%s','%s') %s`, user.PkId, user.UserName, user.IsPrivate, user.FollowFlag, user.CreatedAt, user.ModifiedAt, separateCharset)
		}
	}
	prepare, err := tx.Prepare(fmt.Sprintf(`INSERT INTO users (pk_id,user_name,is_private,follow_flag,created_at,modified_at) VALUES %s;`, insertValue))
	if err != nil {
		return err
	}
	_, err = prepare.Exec()
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func CheckUserExists(store *sql.DB, userPkId string) (bool, error) {
	rows, err := store.Query("SELECT * FROM users WHERE pk_id = ?", userPkId)
	if err != nil {
		return false, err
	}
	var users []UserModel
	for rows.Next() {
		var user UserModel
		err := rows.Scan(&user.PkId, &user.UserName, &user.IsPrivate, &user.FollowFlag, &user.CreatedAt, &user.ModifiedAt)
		if err != nil {
			return false, err
		}
		users = append(users, user)
	}
	return len(users) > 0, err
}
