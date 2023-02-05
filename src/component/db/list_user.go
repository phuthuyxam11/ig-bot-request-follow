package db

import "database/sql"

func GetAllUser(store *sql.DB) ([]UserModel, error) {
	rows, err := store.Query("SELECT * FROM users WHERE follow_flag = 'no'")
	if err != nil {
		return nil, err

	}
	var users []UserModel
	for rows.Next() {
		var user UserModel
		err := rows.Scan(&user.PkId, &user.UserName, &user.IsPrivate, &user.FollowFlag, &user.CreatedAt, &user.ModifiedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}
