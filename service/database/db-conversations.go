package database

// this function checks if the tables for conversations exist. if not, it creates them.
//
// This includes :
//
// conversations_table:
// conversation_id (primary key, INTEGER or UUID)
// is_group (BOOLEAN)
// name (TEXT)
// picture (TEXT)
//
// participants_table:
// participant_id (primary key, INTEGER or UUID)
// conversation_id (foreign key, INTEGER or UUID)
// user_id (foreign key, INTEGER or UUID)
//
// messages_table:
// message_id (primary key, INTEGER or UUID)
// conversation_id (foreign key, INTEGER or UUID)
// sender_id (foreign key, INTEGER or UUID)
// content (TEXT)
// is_forwarded (BOOLEAN)
// timestamp (DATETIME or INTEGER for Unix timestamp)
func (db *appdbimpl) CheckConversationsTable() error {
	// first, check if the conversations_table exists
	var tableName string
	err := db.c.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversations_table';`).Scan(&tableName)
	if err != nil {
		// if the table doesn't exist, create it.
		sqlStmt := `CREATE TABLE conversations_table (conversation_id TEXT PRIMARY KEY, is_group BOOLEAN, name TEXT, picture TEXT);`
		_, err = db.c.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	// next, check if the participants_table exists
	err = db.c.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='participants_table';`).Scan(&tableName)
	if err != nil {
		// if the table doesn't exist, create it.
		sqlStmt := `CREATE TABLE participants_table (participant_id TEXT PRIMARY KEY, conversation_id TEXT, user_id TEXT);`
		_, err = db.c.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	// finally, check if the messages_table exists
	err = db.c.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages_table';`).Scan(&tableName)
	if err != nil {
		// if the table doesn't exist, create it.
		sqlStmt := `CREATE TABLE messages_table (message_id TEXT PRIMARY KEY, conversation_id TEXT, sender_id TEXT, content TEXT, is_forwarded BOOLEAN, timestamp INTEGER);`
		_, err = db.c.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) CreateConversation(conversationID string, isGroup bool, name string, picture string) error {
	db.CheckConversationsTable()
	// if some of the values are nil, they will be set to their default values.
	_, err := db.c.Exec("INSERT INTO conversations_table (conversation_id, is_group, name, picture) VALUES (?, ?, ?, ?)", conversationID, isGroup, name, picture)
	return err
}

// this function updates the conversation with the given conversationID.
// the name and picture are updated to the given values.
func (db *appdbimpl) UpdateConversation(conversationID string, isGroup bool, name string, picture string) error {
	db.CheckConversationsTable()
	// if the function is given nil for some of the values, it will keep the old values.
	_, err := db.c.Exec("UPDATE conversations_table SET is_group = COALESCE(?, is_group), name = COALESCE(?, name), picture = COALESCE(?, picture) WHERE conversation_id = ?", isGroup, name, picture, conversationID)
	return err
}

func (db *appdbimpl) DeleteConversation(conversationID string) error {
	db.CheckConversationsTable()
	_, err := db.c.Exec("DELETE FROM conversations_table WHERE conversation_id = ?", conversationID)
	return err
}

func (db *appdbimpl) GetConversationInfo(conversationID string) (string, bool, string, string, error) {
	db.CheckConversationsTable()
	var isGroup bool
	var name string
	var picture string
	err := db.c.QueryRow("SELECT is_group, name, picture FROM conversations_table WHERE conversation_id = ?", conversationID).Scan(&isGroup, &name, &picture)
	return conversationID, isGroup, name, picture, err
}

func (db *appdbimpl) GetUserConversations(userID string) ([]string, error) {
	db.CheckConversationsTable()
	rows, err := db.c.Query("SELECT conversation_id FROM participants_table WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []string
	for rows.Next() {
		var conversationID string
		err = rows.Scan(&conversationID)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversationID)
	}
	return conversations, nil
}

// add participant to conversation
func (db *appdbimpl) AddParticipant(conversationID string, userID string) error {
	db.CheckConversationsTable()
	_, err := db.c.Exec("INSERT INTO participants_table (participant_id, conversation_id, user_id) VALUES (?, ?, ?)", conversationID+userID, conversationID, userID)
	return err
}

// remove participant from conversation
func (db *appdbimpl) RemoveParticipant(conversationID string, userID string) error {
	db.CheckConversationsTable()
	_, err := db.c.Exec("DELETE FROM participants_table WHERE conversation_id = ? AND user_id = ?", conversationID, userID)
	return err
}

// get all participants in a conversation
func (db *appdbimpl) GetParticipants(conversationID string) ([]string, error) {
	db.CheckConversationsTable()
	rows, err := db.c.Query("SELECT user_id FROM participants_table WHERE conversation_id = ?", conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []string
	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}
	return participants, nil
}
