package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Bot handles communication to and from Messenger
type Bot struct {
	pageAccessToken   string
	verificationToken string
	messagingHandler  func(m Messaging) error
	httpClient        http.Client
}

// New creates a new Bot with the appropriate tokens and messaging handler.
func New(pageAccessToken, verificationToken string, messagingHandler func(m Messaging) error) *Bot {
	return &Bot{
		pageAccessToken:   pageAccessToken,
		verificationToken: verificationToken,
		messagingHandler:  messagingHandler,
	}
}

// HandleWebhookPost executes the messagingHandler callback for each Messaging present.
func (m *Bot) HandleWebhookPost(w http.ResponseWriter, r *http.Request) {
	var webhookData WebhookRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&webhookData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	for _, entry := range webhookData.Entries {
		for _, messaging := range entry.Messagings {
			if messaging.Message != nil {
				go func() {
					err := m.messagingHandler(messaging)
					if err != nil {
						log.Println("Error executing messaging handler:", err)
					}
				}()
			}
		}
	}

	w.WriteHeader(200)
}

// SendGenericTemplateReply uses the Send API to send a message to a user
func (m *Bot) SendGenericTemplateReply(recipientID int64, elements []GenericTemplateElement) (err error) {
	var r GenericTemplateReply

	r.Recipient.ID = recipientID
	r.Message.Attachment.Type = "template"
	r.Message.Attachment.Payload.TemplateType = "generic"
	r.Message.Attachment.Payload.Elements = elements

	marshalled, err := json.Marshal(r)
	if err != nil {
		log.Println("Error marshalling send request:", err)
		return
	}

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages", bytes.NewBuffer(marshalled))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("access_token", m.pageAccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Println("Error sending message request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return
	}

	var sendError SendError

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&sendError)
	if err != nil {
		return
	}

	return fmt.Errorf("Error sending response: %s", sendError.Error.Message)
}

// HandleVerificationChallenge allows Facebook to verify this bot.
func (m *Bot) HandleVerificationChallenge(w http.ResponseWriter, r *http.Request) {
	givenVerificationToken := r.URL.Query().Get("hub.verify_token")

	if givenVerificationToken != m.verificationToken {
		http.Error(w, "Incorrect verification token", http.StatusUnauthorized)
	} else {
		w.Write([]byte(r.URL.Query().Get("hub.challenge")))
	}
}
