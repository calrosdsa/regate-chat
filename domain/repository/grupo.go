package repository

import "context"

type GrupoRepository interface {
	SaveGrupoMessage(ctx context.Context, d *Message) (err error)
	GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []Message, err error)
	GetChatUnreadMessage(ctx context.Context, chatId int, lastUpdate string) (res []Message, err error)
	UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error)
	DeleteMessage(ctx context.Context,id int)(err error)

	GetUsers(ctx context.Context,d RequestUsersGroupOrRoom)(actives []UsersGroupOrRoom,inactives []UsersGroupOrRoom,err error)
}

type GrupoUseCase interface {
	SaveGrupoMessage(ctx context.Context, d *Message) (err error)
	GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []Message,
		nextPage int16, err error)
	GetChatUnreadMessage(ctx context.Context, chatId int, lastUpdate string) (res []Message, err error)
	UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error)
	DeleteMessage(ctx context.Context,id int)(err error)

	GetUsers(ctx context.Context,d RequestUsersGroupOrRoom)(actives []UsersGroupOrRoom,inactives []UsersGroupOrRoom,err error)
}

type Message struct {
	Id          int              `json:"id"`
	LocalId     int64            `json:"local_id,omitempty"`
	ChatId      int              `json:"chat_id,omitempty"`
	ProfileId   int              `json:"profile_id,omitempty"`
	TypeMessage GrupoMessageType `json:"type_message,omitempty"`
	Content     string           `json:"content"`
	Data        *string          `json:"data,omitempty"`
	CreatedAt   string           `json:"created_at,omitempty"`
	ParentId    int              `json:"parent_id,omitempty"`
	ReplyTo     *int             `json:"reply_to,omitempty"`
	IsDeleted   bool             `json:"is_deleted"`
	//Only fon conversation message
	IsUser bool `json:"is_user"`
	IsRead bool `json:"is_read"`
}

type MessageWithReply struct {
	Id          int              `json:"id"`
	LocalId     int64            `json:"local_id,omitempty"`
	ChatId      int              `json:"chat_id,omitempty"`
	ProfileId   int              `json:"profile_id,omitempty"`
	TypeMessage GrupoMessageType `json:"type_message,omitempty"`
	Content     string           `json:"content"`
	Data        *string          `json:"data,omitempty"`
	CreatedAt   string           `json:"created_at,omitempty"`
	ParentId    int              `json:"parent_id,omitempty"`
	ReplyTo     *int             `json:"reply_to,omitempty"`
	IsDeleted   bool             `json:"is_deleted"`
	//Only fon conversation message
	IsUser bool `json:"is_user"`
	IsRead bool `json:"is_read"`
	ReplyMessage ReplyMessage     `json:"reply,omitempty"`
}

type ReplyMessage struct {
	Id          int              `json:"id"`
	Content     string           `json:"content"`
	Data        *string          `json:"data"`
	CreatedAt   string           `json:"created_at"`
	TypeMessage GrupoMessageType `json:"type_message"`
}

const (
	GrupoEventSala    = "sala"
	GrupoEventMessage = "message"
	GrupoEventIgnore  = "ignore"
)

type GrupoMessageType int8

const (
	TypeMessageCommon      = 0
	TypeMessageInstalacion = 1
	TypeMessageSala        = 2
)
