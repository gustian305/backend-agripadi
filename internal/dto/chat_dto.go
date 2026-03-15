package dto

type SendChatRequest struct {
	ConversationID string `json:"conversation_id"`
	ChatType       string `json:"chat_type"` // text / image
	Content        string `json:"content"`
}

type ChatAttachmentResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

type ChatResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderType     string `json:"sender_type"`
	ChatType       string `json:"chat_type"`
	Content        string `json:"content"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`

	AIResult     string  `json:"ai_result,omitempty"`
	AIConfidence float64 `json:"ai_confidence,omitempty"`

	Attachment []ChatAttachmentResponse `json:"attachment,omitempty"`
}
