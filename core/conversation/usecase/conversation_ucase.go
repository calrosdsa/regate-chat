package usecase

import (
	"context"
	"encoding/json"
	"log"
	r "message/domain/repository"
	"time"

	"github.com/segmentio/kafka-go"
)

type grupoUcase struct {
	timeout          time.Duration
	conversationRepo r.ConversationRepository
	utilU             r.UtilUseCase
	kafkaW           *kafka.Writer
}

func NewUseCase(timeout time.Duration, conversationRepo r.ConversationRepository,utilU r.UtilUseCase) r.ConversationUseCase {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9094"),
		Topic:    "notification-message-group",
		Balancer: &kafka.LeastBytes{},
	}
	return &grupoUcase{
		timeout:          timeout,
		conversationRepo: conversationRepo,
		kafkaW:           w,
		utilU: utilU,
	}
}
func (u *grupoUcase) GetOrCreateConversation(ctx context.Context,id int,profileId int)(conversationId int,err error){
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	conversationId,err = u.conversationRepo.GetOrCreateConversation(ctx,id,profileId)
	return 
}

func (u *grupoUcase) GetConversations(ctx context.Context, id int) (res []r.Conversation, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res,err = u.conversationRepo.GetConversations(ctx,id)
	return
}

func (u *grupoUcase) GetMessages(ctx context.Context, id int,page int16,size int8) (res []r.Inbox,nextPage int16, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	page = u.utilU.PaginationValues(page)
	res, err = u.conversationRepo.GetMessages(ctx, id,page,size)
	nextPage = u.utilU.GetNextPage(int8(len(res)),int8(size),page + 1)
	return
}

func (u *grupoUcase) SaveMessage(ctx context.Context, d *r.MessageGrupo) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	err = u.conversationRepo.SaveMessage(ctx, d)
	go func() {
		json, err := json.Marshal(d)
		if err != nil {
			log.Println("Fail to parse", err)
		}
		err = u.kafkaW.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Message"),
				Value: json,
			},
		)
		if err != nil {
			log.Println("failed to write messages:", err)
		}
	}()

	return

}
