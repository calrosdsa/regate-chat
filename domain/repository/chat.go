package repository

import "context"

type ChatUseCase interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat,
		nextPage int16,err error)
}

type ChatRepository interface {
	GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []Chat,err error)
}

type Chat struct {
	Id                 int     `json:"id"`
	Photo              *string `json:"photo"`
	Name               string  `json:"name"`
	LastMessage        *string `json:"last_message,omitempty"`
	LastMessageCreated *string `json:"last_message_created,omitempty"`
	MessagesCount      int     `json:"messages_count,omitempty"`
}
