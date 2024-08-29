package db

func (s Storage) StoreMessageID(chatID, entityID int64, messageID int) error {
	q := `
		INSERT INTO bot_messages (chat_id, message_id, entity_id)
		VALUES (?, ?, ?)
	`

	_, err := s.db.Exec(q, chatID, messageID, entityID)

	return err
}

func (s Storage) GetLastMessageID(chatID, entityID int64) (*int64, error) {
	q := `
		SELECT m.message_id
		FROM bot_messages m
		WHERE m.chat_id = ? AND m.entity_id = ?
		ORDER BY m.sent_at DESC
		LIMIT 1
	`

	var messageID int64
	err := s.db.QueryRow(q, chatID, entityID).Scan(&messageID)

	if IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &messageID, err
}
