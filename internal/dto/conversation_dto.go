package dto

type CreateConversationRequest struct {
	IsBotChat bool `json:"is_bot_chat"`
}

type DeleteConversationRequest struct {
	ConversationID []string `json:"conversation_id"`
}

type ConversationResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	IsBotChat bool   `json:"is_bot_chat"`
	AIState   string `json:"ai_state"`
	CreatedAt string `json:"created_at"`
}

type ConversationListResponse struct {
	ID           string `json:"id"`
	LastMessage  string `json:"last_message"`
	ChatCount int    `json:"chat_count"`
	CreatedAt    string `json:"created_at"`
}

type ConversationDetailResponse struct {
	ID    string         `json:"id"`
	Chats []ChatResponse `json:"chats"`
}