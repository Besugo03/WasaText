package database

// checks the comments table and creates it if it doesn't exist.
// the comments table is a table that stores comments.
// has the following fields:
// comment_id (primary key, INTEGER autoincrement)
// message_id (foreign key, INTEGER)
// sender_id (foreign key, TEXT) (the user_id of the sender)
// content (TEXT)
// timestamp (DATETIME)
func (db *appdbimpl) CheckCommentsTable() error {
	_, err := db.c.Exec(`CREATE TABLE IF NOT EXISTS comments_table (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_id INTEGER,
		sender_id TEXT,
		content TEXT,
		timestamp TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (message_id) REFERENCES messages_table(message_id),
		FOREIGN KEY (sender_id) REFERENCES users_table(user_id)
	)`)
	return err
}

// adds a comment to the comments table.
// the message_id is the message that the comment is attached to.
// the sender_id is the user_id of the sender.
// the content is the content of the comment.
func (db *appdbimpl) AddComment(messageID string, senderID string, content string) error {
	err := db.CheckCommentsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("INSERT INTO comments_table (message_id, sender_id, content) VALUES (?, ?, ?)", messageID, senderID, content)
	return err
}

// gets all the comments for a message.
// returns the comment_ids, message_ids, sender_ids, contents, and timestamps.
func (db *appdbimpl) GetComments(messageID string) (commentIDs []string, messageIDs []string, senderIDs []string, contents []string, timestamps []string, err error) {
	err = db.CheckCommentsTable()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	rows, err := db.c.Query("SELECT comment_id, message_id, sender_id, content, timestamp FROM comments_table WHERE message_id = ?", messageID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentID, messageID, senderID, content, timestamp string
		if err := rows.Scan(&commentID, &messageID, &senderID, &content, &timestamp); err != nil {
			return nil, nil, nil, nil, nil, err
		}
		commentIDs = append(commentIDs, commentID)
		messageIDs = append(messageIDs, messageID)
		senderIDs = append(senderIDs, senderID)
		contents = append(contents, content)
		timestamps = append(timestamps, timestamp)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return commentIDs, messageIDs, senderIDs, contents, timestamps, nil
}

// deletes a comment from the comments table.
func (db *appdbimpl) DeleteComment(commentID string) error {
	err := db.CheckCommentsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM comments_table WHERE comment_id = ?", commentID)
	return err
}

// gets the information for a comment.
// returns the message_id, sender_id, content, and timestamp.
func (db *appdbimpl) GetCommentInfo(commentID string) (string, string, string, string, error) {
	err := db.CheckCommentsTable()
	if err != nil {
		return "", "", "", "", err
	}
	var messageID, senderID, content, timestamp string
	err = db.c.QueryRow("SELECT message_id, sender_id, content, timestamp FROM comments_table WHERE comment_id = ?", commentID).Scan(&messageID, &senderID, &content, &timestamp)
	return messageID, senderID, content, timestamp, err
}
