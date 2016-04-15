package messenger

// GenericTemplateReply represents a generic template response
type GenericTemplateReply struct {
	Recipient Participant `json:"recipient"`
	Message   struct {
		Attachment struct {
			Type    string `json:"type"`
			Payload struct {
				TemplateType string                   `json:"template_type"`
				Elements     []GenericTemplateElement `json:"elements"`
			} `json:"payload"`
		} `json:"attachment"`
	} `json:"message"`
}

// GenericTemplateElement represents a template item in a generic template reply.
type GenericTemplateElement struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ItemURL  string `json:"item_url"`
}

// SendError represents an error returned by the Send API
type SendError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
