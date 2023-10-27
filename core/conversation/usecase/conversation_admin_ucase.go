package usecase

import (
	"context"
	r "message/domain/repository"
	"time"

	"github.com/segmentio/kafka-go"
)

type conversationAdminUseCase struct {
	timeout               time.Duration
	conversationAdminRepo r.ConversationAdminRepository
	utilU                 r.UtilUseCase
	kafkaW                *kafka.Writer
}

func NewAdminUseCase(timeout time.Duration, conversationAdminRepo r.ConversationAdminRepository, utilU r.UtilUseCase) r.ConversationAdminUseCase {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9094"),
		Topic:    "notification-message-group",
		Balancer: &kafka.LeastBytes{},
	}
	return &conversationAdminUseCase{
		timeout:               timeout,
		conversationAdminRepo: conversationAdminRepo,
		kafkaW:                w,
		utilU:                 utilU,
	}
}
func (u *conversationAdminUseCase) GetConversationsEstablecimiento(ctx context.Context,uuid string) (res []r.ChatEstablecimiento, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err = u.conversationAdminRepo.GetConversationsEstablecimiento(ctx, uuid)
	return
}

func (u *conversationAdminUseCase) GetMessages(ctx context.Context, id int, page int16,
	size int8) (res []r.Message, nextPage int16, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	page = u.utilU.PaginationValues(page)
	res, err = u.conversationAdminRepo.GetMessages(ctx, id, page, int8(size))
	if err != nil {
		u.utilU.LogError("SaveGrupoMessage", "grupo_usecase", err.Error())
	}
	nextPage = u.utilU.GetNextPage(int8(len(res)), int8(size), page+1)
	return
}
