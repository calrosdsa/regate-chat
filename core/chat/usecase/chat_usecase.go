package usecase

import (
	"context"
	"encoding/json"
	"log"
	r "message/domain/repository"
	ws "message/domain/ws"
	"time"
	// "github.com/segmentio/kafka-go"
)

type chatUseCase struct {
	timeout  time.Duration
	chatRepo r.ChatRepository
	utilU    r.UtilUseCase
	grupoU   r.GrupoUseCase
	wsServer *ws.WsServer

	// kafkaW           *kafka.Writer
}

func NewUseCase(timeout time.Duration, charRepo r.ChatRepository, utilU r.UtilUseCase,
	grupoU r.GrupoUseCase, wsServer *ws.WsServer) r.ChatUseCase {
	// w := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9094"),
	// 	Topic:    "notification-message-group",
	// 	Balancer: &kafka.LeastBytes{},
	// }
	return &chatUseCase{
		timeout:  timeout,
		chatRepo: charRepo,
		wsServer: wsServer,
		grupoU:   grupoU,
		// kafkaW:           w,
		utilU: utilU,
	}
}

func (u *chatUseCase) PublishMessage(ctx context.Context, msg r.MessagePublishRequest) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	log.Println(msg.TypeChat)
	switch msg.TypeChat {
	
	case r.TypeChatGrupo:
			err := u.grupoU.SaveGrupoMessage(ctx, &msg.Message)
			if err != nil {
				u.utilU.LogError("PublisMessage_SaveGrupoMessage", "chat_usecase", err.Error())
				return
			}
			messagePaylaod, err := json.Marshal(msg.Message)
			if err != nil {
				u.utilU.LogError("PublisMessage_Marshal", "chat_usecase", err.Error())
				return
			}
			event := r.MessageEvent{
				Type:    "message",
				Payload: string(messagePaylaod),
			}
			payload, err := json.Marshal(event)
			if err != nil {
				u.utilU.LogError("PublisMessage_Marshal", "chat_usecase", err.Error())
				return
			}
			u.wsServer.Publish(payload, msg.ChatId)
	}
}

func (u *chatUseCase) GetChatsUser(ctx context.Context, profileId int, page int16, size int8) (res []r.Chat,
	nextPage int16, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	page = u.utilU.PaginationValues(page)
	res, err = u.chatRepo.GetChatsUser(ctx, profileId, page, int8(size))
	if err != nil {
		u.utilU.LogError("GetChatUser", "chat_usecase", err.Error())
	}
	nextPage = u.utilU.GetNextPage(int8(len(res)), int8(size), page+1)
	return
}
