package database

import "strconv"

// this function checks if the tables for conversations exist. if not, it creates them.
//
// This includes :
//
// conversations_table:
// conversation_id (primary key, INTEGER or UUID)
// is_group (BOOLEAN)
// name (TEXT)
// photo (TEXT)
//
// participants_table:
// participant_id (the user_id of the participant to the conversation) (primary key, INTEGER or UUID)
// conversation_id (foreign key, INTEGER or UUID)
//
// messages_table:
// message_id (primary key, INTEGER or UUID)
// conversation_id (foreign key, INTEGER or UUID)
// sender_id (foreign key, INTEGER or UUID)
// content (TEXT)
// is_forwarded (BOOLEAN)
// timestamp (DATETIME or INTEGER for Unix timestamp)
func (db *appdbimpl) CheckConversationsTable() error {
	_, err := db.c.Exec(`CREATE TABLE IF NOT EXISTS conversations_table (
        conversation_id INTEGER PRIMARY KEY AUTOINCREMENT,
        is_group BOOLEAN,
        name TEXT,
        photo TEXT
    )`)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`CREATE TABLE IF NOT EXISTS participants_table (
		conversation_id INTEGER,
		user_id TEXT,
		FOREIGN KEY (conversation_id) REFERENCES conversations_table(conversation_id),
		FOREIGN KEY (user_id) REFERENCES users_table(user_id)
	)`)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`CREATE TABLE IF NOT EXISTS messages_table (
        message_id INTEGER PRIMARY KEY AUTOINCREMENT,
        conversation_id INTEGER,
        sender_id TEXT,
        content TEXT,
        timestamp TEXT DEFAULT CURRENT_TIMESTAMP,
		is_forwarded BOOLEAN,
        FOREIGN KEY (conversation_id) REFERENCES conversations_table(conversation_id),
        FOREIGN KEY (sender_id) REFERENCES users_table(user_id)
    )`)
	return err
}

// returns the newly created conversation's conversation_id.
func (db *appdbimpl) CreateConversation(isGroup bool, name string, photo string) (string, error) {
	db.CheckConversationsTable()
	// if some of the values are nil, they will be set to their default values.
	res, err := db.c.Exec("INSERT INTO conversations_table (is_group, name, photo) VALUES (?, ?, ?)", isGroup, name, photo)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	return strconv.FormatInt(id, 10), err
}

// this function updates the conversation with the given conversationID.
// the name and photo are updated to the given values.
func (db *appdbimpl) UpdateConversation(conversationID string, isGroup bool, name string, photo string) error {
	db.CheckConversationsTable()
	// Get the current values
	var currentName, currentPhoto string
	err := db.c.QueryRow("SELECT name, photo FROM conversations_table WHERE conversation_id = ?", conversationID).Scan(&currentName, &currentPhoto)
	if err != nil {
		return err
	}

	// Use the current values if the input is empty
	if name == "" {
		name = currentName
	}
	if photo == "" {
		photo = currentPhoto
	}

	_, err = db.c.Exec("UPDATE conversations_table SET is_group = ?, name = ?, photo = ? WHERE conversation_id = ?", isGroup, name, photo, conversationID)
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
	var photo string
	err := db.c.QueryRow("SELECT is_group, name, photo FROM conversations_table WHERE conversation_id = ?", conversationID).Scan(&isGroup, &name, &photo)
	return conversationID, isGroup, name, photo, err
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
	_, err := db.c.Exec("INSERT INTO participants_table (conversation_id, user_id) VALUES (?, ?)", conversationID, userID)
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

// sends a message to a conversation. Returns the message_id of the newly created message.
func (db *appdbimpl) SendMessage(conversationID string, senderID string, content string, is_forwarded bool) (string, error) {
	db.CheckConversationsTable()
	res, err := db.c.Exec("INSERT INTO messages_table (conversation_id, sender_id, content, is_forwarded) VALUES (?, ?, ?, ?)", conversationID, senderID, content, is_forwarded)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	return strconv.FormatInt(id, 10), err
}

// get information about the messages in a conversation
func (db *appdbimpl) GetMessages(conversationID string) ([]string, []string, []string, []string, []string, error) {
	db.CheckConversationsTable()
	rows, err := db.c.Query("SELECT message_id, conversation_id, sender_id, content, timestamp FROM messages_table WHERE conversation_id = ?", conversationID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer rows.Close()

	var messageIDs []string
	var conversationIDs []string
	var senderIDs []string
	var contents []string
	var timestamps []string
	for rows.Next() {
		var messageID string
		var conversationID string
		var senderID string
		var content string
		var timestamp string
		err = rows.Scan(&messageID, &conversationID, &senderID, &content, &timestamp)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		messageIDs = append(messageIDs, messageID)
		conversationIDs = append(conversationIDs, conversationID)
		senderIDs = append(senderIDs, senderID)
		contents = append(contents, content)
		timestamps = append(timestamps, timestamp)
	}
	return messageIDs, conversationIDs, senderIDs, contents, timestamps, err
}

// delete a message from a conversation
func (db *appdbimpl) DeleteMessage(messageID string) error {
	db.CheckConversationsTable()
	_, err := db.c.Exec("DELETE FROM messages_table WHERE message_id = ?", messageID)
	return err
}

// get information about a message
func (db *appdbimpl) GetMessageInfo(messageID string) (string, string, string, string, error) {
	db.CheckConversationsTable()
	var conversationID string
	var senderID string
	var content string
	var timestamp string
	err := db.c.QueryRow("SELECT conversation_id, sender_id, content, timestamp FROM messages_table WHERE message_id = ?", messageID).Scan(&conversationID, &senderID, &content, &timestamp)
	return conversationID, senderID, content, timestamp, err
}

// get messages from a conversation
func (db *appdbimpl) GetConversationMessages(conversationID string) ([]string, []string, []string, []string, []string, error) {
	db.CheckConversationsTable()
	rows, err := db.c.Query("SELECT message_id, conversation_id, sender_id, content, timestamp FROM messages_table WHERE conversation_id = ?", conversationID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer rows.Close()

	var messageIDs []string
	var conversationIDs []string
	var senderIDs []string
	var contents []string
	var timestamps []string
	for rows.Next() {
		var messageID string
		var conversationID string
		var senderID string
		var content string
		var timestamp string
		err = rows.Scan(&messageID, &conversationID, &senderID, &content, &timestamp)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		messageIDs = append(messageIDs, messageID)
		conversationIDs = append(conversationIDs, conversationID)
		senderIDs = append(senderIDs, senderID)
		contents = append(contents, content)
		timestamps = append(timestamps, timestamp)
	}
	return messageIDs, conversationIDs, senderIDs, contents, timestamps, err
}
