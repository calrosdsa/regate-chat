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
	ChatId         int      `json:"chat_id"`
	LastUpdateChat string   `json:"last_update_chat"`
	TypeChat       TypeChat `json:"type_chat"`
}
type ChatUseCase interface {
	GetChatByParentId(ctx context.Context,parentId int,typeChat TypeChat)(res Chat,err error)
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat,
		nextPage int16, err error)
	PublishMessage(ctx context.Context, msg MessagePublishRequest) (res int, err error)
	GetChatUnreadMessages(ctx context.Context, d RequestChatUnreadMessages) (res []Message, err error)

	DeleteMessage(ctx context.Context, d DeleteMessageRequet) (err error)
	GetDeletedMessages(ctx context.Context, id int) (res []int, err error)
	GetUsers(ctx context.Context,d RequestUsersGroupOrRoom)(res []UsersGroupOrRoom,err error)
	NotifyNewUser(chatId int,d UsersGroupOrRoom)
}

type ChatRepository interface {
	GetChatByParentId(ctx context.Context,parentId int,typeChat TypeChat)(res Chat,err error)
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat, err error)
	DeleteMessage(ctx context.Context, id int, chatId int) (err error)
	GetDeletedMessages(ctx context.Context, id int) (res []int, err error)
}

type RequestUsersGroupOrRoom struct {
	ParentId int      `json:"parent_id"`
	TypeChat TypeChat `json:"type_chat"`
	//cantidad de usuarios almacenadosasa localmente
	LastUpdated *string `json:"last_updated"`
}

type DeleteMessageRequet struct {
	Id       int      `json:"id"`
	Ids      []int    `json:"ids"`
	ChatId   int      `json:"chat_id"`
	TypeChat TypeChat `json:"type_chat"`
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

type WskafkaPayload struct {
	Payload string             `json:"payload"`
	Type    WskafkaPayloadType `json:"type"`
}
type WskafkaPayloadType int8

const (
	WsKafkaTypeDeleteMessage = 1
)

type MessageEventType string

const (
	EventTypeMessage        MessageEventType = "message"
	EventTypeDeletedMessage MessageEventType = "delete-message"
)

type TypeChat int8

const (
	TypeChatGrupo                = 1
	TypeChatInboxEstablecimiento = 2
	TypeChatSala                 = 3
)
