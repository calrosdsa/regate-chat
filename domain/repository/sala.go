package repository

import "context"

type SalaUseCase interface {
	SaveMessage(ctx context.Context, d *Message) (err error)
	GetChatUnreadMessages(ctx context.Context, chatId int, lastUpdate string) (res []Message, err error)
	DeleteMessage(ctx context.Context,id int)(err error)
}

type SalaRepository interface {
	SaveMessage(ctx context.Context, d *Message) (err error)
	// GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []MessageGrupo, err error)
	GetChatUnreadMessages(ctx context.Context, chatId int, lastUpdate string) (res []Message, err error)
	DeleteMessage(ctx context.Context,id int)(err error)
}
