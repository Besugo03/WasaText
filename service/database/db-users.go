package database

import "errors"

// function that makes sure the database exists. if not, it creates the database.
func (db *appdbimpl) CheckUsersTable() error {
	// the table is called users_table.
	// you CANNOT have two users with the same username.
	// you CANNOT have two users with the same user_id.
	// the fields are :
	// username (TEXT)
	// user_id (TEXT)
	// photo (TEXT)

	// if the table doesn't exist, create it.
	_, err := db.c.Exec(`CREATE TABLE IF NOT EXISTS users_table (
		username TEXT PRIMARY KEY,
		user_id TEXT,
		photo TEXT
	)`)
	return err
}

// this basically inserts a new user into the database.
// the username is the user's username and the userid is the user's auth token.
func (db *appdbimpl) RegisterUser(username string, userid string, photo string) error {
	db.CheckUsersTable()
	// _, err := db.c.Exec("INSERT INTO users_table (username, user_id) VALUES (?, ?)", username, userid)
	// if we try to insert a user that already exists, the database will return an error.
	// if we try to insert a userid that already exists, the database will return an error.

	// first, check if the user or userid already exists
	// if it does, return an error.
	// if it doesn't, insert the user into the database.
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ? OR user_id = ?", username, userid).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("user already exists")
	}
	// photo is optional. if it's not provided, it will be set to an empty string.
	_, err = db.c.Exec("INSERT INTO users_table (username, user_id, photo) VALUES (?, ?, ?)", username, userid, photo)
	return err
}

func (db *appdbimpl) GetUserID(username string) (string, error) {
	db.CheckUsersTable()
	var userid string
	// get the userid from the table
	err := db.c.QueryRow("SELECT user_id FROM users_table WHERE username = ?", username).Scan(&userid)
	// this basically gets the user's auth token from the database.
	// the username is the user's username.
	// the userid is the user's auth token.
	// asking the database for a username that is not in the database will return an error.
	return userid, err
}

func (db *appdbimpl) GetUserName(userid string) (string, error) {
	db.CheckUsersTable()
	var name string
	err := db.c.QueryRow("SELECT username FROM users_table WHERE user_id = ?", userid).Scan(&name)
	return name, err
}

func (db *appdbimpl) UpdateUser(username string, userid string, photo string) error {
	db.CheckUsersTable()
	// if the user doesn't exist, return an error.
	// if the user exists, update the user.
	// for each argument that is nil, keep the old value.
	// first, we get the old values.
	var oldUsername string
	var oldPhoto string
	err := db.c.QueryRow("SELECT username, photo FROM users_table WHERE user_id = ?", userid).Scan(&oldUsername, &oldPhoto)
	if err != nil {
		return err
	}
	// if the username is nil or "", keep the old value.
	if username == "" {
		username = oldUsername
	}
	// if the photo is nil or "", keep the old value.
	if photo == "" {
		photo = oldPhoto
	}
	// update the user in the database.
	_, err = db.c.Exec("UPDATE users_table SET username = ?, photo = ? WHERE user_id = ?", username, photo, userid)
	return err
}

// get all users
func (db *appdbimpl) GetAllUsers() ([]string, error) {
	db.CheckUsersTable()
	rows, err := db.c.Query("SELECT username FROM users_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		users = append(users, username)
	}
	return users, nil
}

// check if a user exists
func (db *appdbimpl) CheckUserExistence(username string) (bool, error) {
	db.CheckUsersTable()
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
