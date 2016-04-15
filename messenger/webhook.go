package messenger

// WebhookRequest represents a webhook sent to the bot from Messenger
type WebhookRequest struct {
	Entries []Entry `json:"entry"`
}

// Entry represents a batch of Messagings
type Entry struct {
	ID         int64       `json:"id"`
	Time       int64       `json:"time"`
	Messagings []Messaging `json:"messaging"`
}

// Messaging reprents a single message to the bot
type Messaging struct {
	Sender    Participant `json:"sender"`
	Recipient Participant `json:"recipient"`
	Timestamp int64       `json:"timestamp"`
	Message   *Message    `json:"message,omitempty"`
}

// Participant contains the ID of a sender or recipient of a message
type Participant struct {
	ID int64 `json:"id"`
}

// Message contains the contents of a message sent to the bot (text only for now)
type Message struct {
	ID             string `json:"mid"`
	SequenceNumber int64  `json:"seq"`
	Text           string `json:"text"`
}
