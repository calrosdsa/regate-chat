package repository

import "context"

type SalaUseCase interface {
	SaveMessage(ctx context.Context, d *MessageGrupo) (err error)
	GetChatUnreadMessage(ctx context.Context, chatId int64, lastUpdate string) (res []MessageGrupo, err error)
}

type SalaRepository interface {
	SaveMessage(ctx context.Context, d *MessageGrupo) (err error)
	// GetUnreadMessages(ctx context.Context, profileId int, page int16, size int8) (res []MessageGrupo, err error)
	GetChatUnreadMessage(ctx context.Context, chatId int64, lastUpdate string) (res []MessageGrupo, err error)
}