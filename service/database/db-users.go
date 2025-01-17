package database

import "errors"

// function that makes sure the database exists. if not, it creates the database.
func (db *appdbimpl) CheckUsersTable() error {
	// the table is called users_table.
	// you CANNOT have two users with the same username.
	// you CANNOT have two users with the same user_id.
	var tableName string
	err := db.c.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users_table';`).Scan(&tableName)
	if err != nil {
		// if the table doesn't exist, create it.
		sqlStmt := `CREATE TABLE users_table (username TEXT PRIMARY KEY, user_id TEXT);`
		_, err = db.c.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// this basically inserts a new user into the database.
// the username is the user's username and the userid is the user's auth token.
func (db *appdbimpl) SetUser(username, userid string) error {
	db.CheckUsersTable()
	// _, err := db.c.Exec("INSERT INTO users_table (username, user_id) VALUES (?, ?)", username, userid)
	// if we try to insert a user that already exists, the database will return an error.
	// if we try to insert a userid that already exists, the database will return an error.

	// first, check if the user or userid already exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users_table WHERE username = ? OR user_id = ?", username, userid).Scan(&count)
	if count > 0 {
		return errors.New("user already exists")
	}
	if err != nil {
		return err
	}
	// if the user or userid doesn't exist, insert the user into the database.
	_, err = db.c.Exec("INSERT INTO users_table (username, user_id) VALUES (?, ?)", username, userid)
	return err
}

func (db *appdbimpl) GetUser(username string) (string, error) {
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
