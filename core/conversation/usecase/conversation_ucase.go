package usecase

import (
	"context"
	// "encoding/json"
	// "log"
	r "message/domain/repository"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type conversationUCase struct {
	timeout          time.Duration
	conversationRepo r.ConversationRepository
	utilU            r.UtilUseCase
	kafkaW           *kafka.Writer
}

func NewUseCase(timeout time.Duration, conversationRepo r.ConversationRepository, utilU r.UtilUseCase) r.ConversationUseCase {
	w := &kafka.Writer{
		Addr:     kafka.TCP(viper.GetString("kafka.host")),
		Topic:    "notification-message-conversation",
		Balancer: &kafka.LeastBytes{},
	}
	return &conversationUCase{
		timeout:          timeout,
		conversationRepo: conversationRepo,
		kafkaW:           w,
		utilU:            utilU,
	}
}
func (u *conversationUCase) GetOrCreateConversation(ctx context.Context, id int, profileId int) (conversationId int, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	conversationId, err = u.conversationRepo.GetOrCreateConversation(ctx, id, profileId)
	return
}


func (u *conversationUCase) GetConversations(ctx context.Context, id int) (res []r.Conversation, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err = u.conversationRepo.GetConversations(ctx, id)
	if err != nil {
		u.utilU.LogError("GetConversations","conversation_usecase",err.Error())
	}
	return
}

func (u *conversationUCase) GetMessages(ctx context.Context, id int, page int16, size int8) (res []r.Inbox, nextPage int16, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	page = u.utilU.PaginationValues(page)
	res, err = u.conversationRepo.GetMessages(ctx, id, page, size)
	if err != nil {
		u.utilU.LogError("GetMessages","conversation_usecase",err.Error())
	}
	nextPage = u.utilU.GetNextPage(int8(len(res)), int8(size), page+1)
	return
}

func (u *conversationUCase) SaveMessage(ctx context.Context, d *r.Message) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	err = u.conversationRepo.SaveMessage(ctx, d)
	if err != nil {
		u.utilU.LogError("SaveMessage","conversation_usecase",err.Error())
	}
	go u.utilU.SendMessageToKafka(u.kafkaW,d,"Message")
	// go func() {
	// 	json, err := json.Marshal(d)
	// 	if err != nil {
	// 		log.Println("Fail to parse", err)
	// 	}
	// 	err = u.kafkaW.WriteMessages(context.Background(),
	// 		kafka.Message{
	// 			Key:   []byte("Message"),
	// 			Value: json,
	// 		},
	// 	)
	// 	if err != nil {
	// 		log.Println("failed to write messages:", err)
	// 	}
	// }()
	return
}


func (u *conversationUCase)UpdateMessagesToReaded(ctx context.Context,ids []int)(err error){
	ctx,cancel := context.WithTimeout(ctx,u.timeout)
	defer cancel()
	for i :=0;i < len(ids);i++ {
		err = u.conversationRepo.UpdateMessageToReaded(ctx,ids[i])
		if err != nil {
			u.utilU.LogError("UpdateMessagesToReaded","conversation_usecase",err.Error())
		}
	}
	return
}

func (u *conversationUCase)DeleteMessage(ctx context.Context,id int)(err error){
	ctx,cancel := context.WithTimeout(ctx,u.timeout)
	defer cancel()
	err = u.conversationRepo.DeleteMessage(ctx,id)
	return
}

