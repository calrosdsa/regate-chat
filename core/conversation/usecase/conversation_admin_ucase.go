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
func (u *conversationAdminUseCase) GetConversationsEstablecimiento(ctx context.Context,uuid string) (res []r.EstablecimientoConversation, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err = u.conversationAdminRepo.GetConversationsEstablecimiento(ctx, uuid)
	return
}
