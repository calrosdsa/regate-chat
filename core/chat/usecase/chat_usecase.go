package usecase

import (
	"context"
	r "message/domain/repository"
	"time"

	// "github.com/segmentio/kafka-go"
)

type chatUseCase struct {
	timeout  time.Duration
	chatRepo r.ChatRepository
	utilU    r.UtilUseCase
	// kafkaW           *kafka.Writer
}

func NewUseCase(timeout time.Duration,charRepo r.ChatRepository, utilU r.UtilUseCase) r.ChatUseCase {
	// w := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9094"),
	// 	Topic:    "notification-message-group",
	// 	Balancer: &kafka.LeastBytes{},
	// }
	return &chatUseCase{
		timeout:          timeout,
		chatRepo: charRepo,
		// kafkaW:           w,
		utilU:            utilU,
	}
}

func (u *chatUseCase) GetChatsUser(ctx context.Context,profileId int,page int16,size int8)(res []r.Chat,
	nextPage int16,err error){
	ctx,cancel := context.WithTimeout(ctx,u.timeout)
	defer cancel()
	page = u.utilU.PaginationValues(page)
	res, err = u.chatRepo.GetChatsUser(ctx, profileId, page, int8(size))
	if err != nil {
		u.utilU.LogError("GetChatUser","chat_usecase",err.Error())
	}
	nextPage = u.utilU.GetNextPage(int8(len(res)), int8(size), page+1)
	return
}
