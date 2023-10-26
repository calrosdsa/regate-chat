package repository

import "context"

type MessagePublishRequest struct {
	Message  Message  `json:"message"`
	TypeChat TypeChat `json:"type_chat"`
	ChatId   int      `json:"chat_id"`
	// EventType MessageEventType `json:"event_type"`
}
type MessageEvent struct {
	Type    MessageEventType `json:"type"`
	Payload string           `json:"payload"`
	// Sala    SalaData     `json:"sala,omitempty"`
}

type RequestChatUnreadMessages struct {
	ChatId         int64    `json:"chat_id"`
	LastUpdateChat string   `json:"last_update_chat"`
	TypeChat       TypeChat `json:"type_chat"`
}
type ChatUseCase interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat,
		nextPage int16, err error)
	PublishMessage(ctx context.Context, msg MessagePublishRequest) (res int, err error)
	GetChatUnreadMessages(ctx context.Context, d RequestChatUnreadMessages) (res []Message, err error)
}

type ChatRepository interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat, err error)
}

type Chat struct {
	Id                 int     `json:"id"`
	Photo              *string  `json:"photo"`
	Name               string   `json:"name"`
	LastMessage        *string  `json:"last_message,omitempty"`
	LastMessageCreated *string  `json:"last_message_created,omitempty"`
	MessagesCount      int      `json:"messages_count,omitempty"`
	TypeChat           TypeChat `json:"type_chat"`
	ParentId           int      `json:"parent_id"`
}

type MessageEventType string

const (
	EventTypeMessage MessageEventType = "message"
)

type TypeChat int8

const (
	TypeChatGrupo                = 1
	TypeChatInboxEstablecimiento = 2
	TypeChatSala                 = 3
)
