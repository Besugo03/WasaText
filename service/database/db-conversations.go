package database

import (
	"strconv"
)

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
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`CREATE TABLE IF NOT EXISTS message_status_table (
        message_id INTEGER,
        user_id TEXT,
        is_received BOOLEAN DEFAULT FALSE,
        is_read BOOLEAN DEFAULT FALSE,
        FOREIGN KEY (message_id) REFERENCES messages_table(message_id),
        FOREIGN KEY (user_id) REFERENCES users_table(user_id)
    )`)
	return err
}

// returns the newly created conversation's conversation_id.
func (db *appdbimpl) CreateConversation(isGroup bool, name string, photo string) (string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return "", err
	}
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
	err := db.CheckConversationsTable()
	if err != nil {
		return err
	}
	// Get the current values
	var currentName, currentPhoto string
	err = db.c.QueryRow("SELECT name, photo FROM conversations_table WHERE conversation_id = ?", conversationID).Scan(&currentName, &currentPhoto)
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
	err := db.CheckConversationsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM conversations_table WHERE conversation_id = ?", conversationID)
	return err
}

func (db *appdbimpl) GetConversationInfo(conversationID string, userID string) (string, bool, string, string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return "", false, "", "", err
	}
	var isGroup bool
	var name string
	var photo string
	err = db.c.QueryRow("SELECT is_group, name, photo FROM conversations_table WHERE conversation_id = ?", conversationID).Scan(&isGroup, &name, &photo)
	if err != nil {
		return "", false, "", "", err
	}

	// Mark all messages in the conversation as received by the user
	_, err = db.c.Exec("UPDATE message_status_table SET is_received = TRUE WHERE message_id IN (SELECT message_id FROM messages_table WHERE conversation_id = ?) AND user_id = ?", conversationID, userID)
	if err != nil {
		return "", false, "", "", err
	}

	return conversationID, isGroup, name, photo, nil
}

func (db *appdbimpl) GetUserConversations(userID string) ([]string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return nil, err
	}
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
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return conversations, nil
}

// add participant to conversation
func (db *appdbimpl) AddParticipant(conversationID string, userID string) error {
	err := db.CheckConversationsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("INSERT INTO participants_table (conversation_id, user_id) VALUES (?, ?)", conversationID, userID)
	if err != nil {
		return err
	}

	// if adding the participant makes it so they are more than 2, the conversation is automatically a group conversation.
	_, err = db.c.Exec("UPDATE conversations_table SET is_group = TRUE WHERE conversation_id = ? AND (SELECT COUNT(*) FROM participants_table WHERE conversation_id = ?) > 2", conversationID, conversationID)
	return err
}

// remove participant from conversation
func (db *appdbimpl) RemoveParticipant(conversationID string, userID string) error {
	err := db.CheckConversationsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM participants_table WHERE conversation_id = ? AND user_id = ?", conversationID, userID)

	if err != nil {
		return err
	}

	// if removing the participant makes it so there are only 2 participants, the conversation is automatically flagged as a non group
	_, err = db.c.Exec("UPDATE conversations_table SET is_group = FALSE WHERE conversation_id = ? AND (SELECT COUNT(*) FROM participants_table WHERE conversation_id = ?) <= 2", conversationID, conversationID)

	return err
}

// get all participants in a conversation
func (db *appdbimpl) GetParticipants(conversationID string) ([]string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return nil, err
	}
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
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return participants, nil
}

// sends a message to a conversation. Returns the message_id of the newly created message.
func (db *appdbimpl) SendMessage(conversationID string, senderID string, content string, is_forwarded bool) (string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return "", err
	}
	res, err := db.c.Exec("INSERT INTO messages_table (conversation_id, sender_id, content, is_forwarded) VALUES (?, ?, ?, ?)", conversationID, senderID, content, is_forwarded)
	if err != nil {
		return "", err
	}
	messageID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	// Get all participants in the conversation
	participants, err := db.GetParticipants(conversationID)
	if err != nil {
		return "", err
	}

	// Insert entries into message_status_table for each participant
	for _, userID := range participants {
		// if the user is the sender, the message is automatically received and read
		if userID == senderID {
			_, err = db.c.Exec("INSERT INTO message_status_table (message_id, user_id, is_received, is_read) VALUES (?, ?, TRUE, TRUE)", messageID, userID)
			if err != nil {
				return "", err
			}
			continue
		}
		_, err = db.c.Exec("INSERT INTO message_status_table (message_id, user_id) VALUES (?, ?)", messageID, userID)
		if err != nil {
			return "", err
		}
	}

	return strconv.FormatInt(messageID, 10), nil
}

// get information about the messages in a conversation
func (db *appdbimpl) GetMessages(conversationID string) ([]string, []string, []string, []string, []string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
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
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return messageIDs, conversationIDs, senderIDs, contents, timestamps, err
}

// delete a message from a conversation
func (db *appdbimpl) DeleteMessage(messageID string) error {
	err := db.CheckConversationsTable()
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM messages_table WHERE message_id = ?", messageID)
	return err
}

// get information about a message
func (db *appdbimpl) GetMessageInfo(messageID string, userID string) (string, string, string, string, string, bool, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return "", "", "", "", "", false, err
	}
	var is_forwarded bool
	var conversationID, senderID, content, timestamp string
	err = db.c.QueryRow("SELECT conversation_id, sender_id, content, timestamp, is_forwarded FROM messages_table WHERE message_id = ?", messageID).Scan(&conversationID, &senderID, &content, &timestamp, &is_forwarded)
	if err != nil {
		return "", "", "", "", "", false, err
	}

	// Mark the message as read by the user
	_, err = db.c.Exec("UPDATE message_status_table SET is_read = TRUE WHERE message_id = ? AND user_id = ?", messageID, userID)
	if err != nil {
		return "", "", "", "", "", false, err
	}

	return conversationID, messageID, senderID, content, timestamp, is_forwarded, nil
}

// get messages from a conversation
func (db *appdbimpl) GetConversationMessages(conversationID string) ([]string, []string, []string, []string, []string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
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
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return messageIDs, conversationIDs, senderIDs, contents, timestamps, err
}

// get the last message from a conversation, along with the sender and the timestamp.
// this is used for the preview of the conversation.
func (db *appdbimpl) GetLastMessage(conversationID string) (string, string, string, string, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return "", "", "", "", err
	}
	var messageID, senderID, timestamp, content string
	err = db.c.QueryRow("SELECT message_id, sender_id, timestamp, content FROM messages_table WHERE conversation_id = ? ORDER BY timestamp DESC LIMIT 1", conversationID).Scan(&messageID, &senderID, &timestamp, &content)
	return content, messageID, senderID, timestamp, err
}

// get the read/received status of a message by a user
func (db *appdbimpl) GetMessageStatus(messageID string, userID string) (bool, bool, bool, bool, error) {
	err := db.CheckConversationsTable()
	if err != nil {
		return false, false, false, false, err
	}
	var isReceived, isRead bool
	err = db.c.QueryRow("SELECT is_received, is_read FROM message_status_table WHERE message_id = ? AND user_id = ?", messageID, userID).Scan(&isReceived, &isRead)
	if err != nil {
		return false, false, false, false, err
	}

	var isReceivedByAll, isReadByAll bool
	// if a message is received by all participants (for each entry of a given message_id in the message_status_table, is_received is TRUE for that message_id), then isReceivedByAll is TRUE

	// 1. get all the entries for the message_id in the message_status_table
	// 2. check if all the entries have is_received as TRUE
	// 3. if all the entries have is_received as TRUE, then isReceivedByAll is TRUE
	rows, err := db.c.Query("SELECT is_received FROM message_status_table WHERE message_id = ?", messageID)
	if err != nil {
		return false, false, false, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var received bool
		err = rows.Scan(&received)
		if err != nil {
			return false, false, false, false, err
		}
		if !received {
			isReceivedByAll = false
			break
		}
		isReceivedByAll = true
	}
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return false, false, false, false, err
	}

	// if a message is read by all participants (for each entry of a given message_id in the message_status_table, is_read is TRUE for that message_id), then isReadByAll is TRUE
	rows, err = db.c.Query("SELECT is_read FROM message_status_table WHERE message_id = ?", messageID)
	if err != nil {
		return false, false, false, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var read bool
		err = rows.Scan(&read)
		if err != nil {
			return false, false, false, false, err
		}
		// fmt.Println(read)
		if !read {
			isReadByAll = false
			break
		}
		isReadByAll = true
	}
	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return false, false, false, false, err
	}

	return isReceived, isRead, isReceivedByAll, isReadByAll, err
}
