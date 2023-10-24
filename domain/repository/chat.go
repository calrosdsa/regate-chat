package repository

import "context"

type MessagePublishRequest struct {
	Message   MessageGrupo          `json:"message"`
	TypeChat  TypeChat         `json:"type_chat"`
	ChatId    int              `json:"chat_id"`
	// EventType MessageEventType `json:"event_type"`
}
type MessageEvent struct {
	Type    MessageEventType `json:"type"`
	Payload string           `json:"payload"`
	// Sala    SalaData     `json:"sala,omitempty"`
}

type ChatUseCase interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat,
		nextPage int16, err error)
	PublishMessage(ctx context.Context, msg MessagePublishRequest)
}

type ChatRepository interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat, err error)
}

type Chat struct {
	Id                 int      `json:"id"`
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
)
