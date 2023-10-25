package repository

import "context"

type GrupoRepository interface {
	SaveGrupoMessage(ctx context.Context, d *MessageGrupo) (err error)
	GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []MessageGrupo, err error)
	GetChatUnreadMessage(ctx context.Context,chatId int64,lastUpdate string)(res []MessageGrupo,err error)
	UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error)
}

type GrupoUseCase interface {
	SaveGrupoMessage(ctx context.Context, d *MessageGrupo) (err error)
	GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []MessageGrupo,
		nextPage int16, err error)
	GetChatUnreadMessage(ctx context.Context,chatId int64,lastUpdate string)(res []MessageGrupo,err error)
	UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error)
}

type MessageGrupo struct {
	Id          int64            `json:"id"`
	ChatId      int              `json:"chat_id"`
	ProfileId   int              `json:"profile_id"`
	TypeMessage GrupoMessageType `json:"type_message"`
	Content     string           `json:"content"`
	Data        *string          `json:"data"`
	CreatedAt   string           `json:"created_at,omitempty"`
	ParentId    int              `json:"parent_id"`
	ReplyTo     *int             `json:"reply_to"`
	// ReplyMessage ReplyMessage     `json:"reply_message"`
}

type ReplyMessage struct {
	Id          int              `json:"id"`
	GrupoId     int              `json:"grupo_id"`
	ProfileId   int              `json:"profile_id"`
	TypeMessage GrupoMessageType `json:"type_message"`
	Data        *string          `json:"data"`
	Content     string           `json:"content"`
	CreatedAt   string           `json:"created_at"`
	ReplyTo     *int             `json:"reply_to"`
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
