package usecase

import (
	"context"
	// "encoding/json"
	"log"
	r "message/domain/repository"
	ws "message/domain/ws"
	"time"

	"github.com/goccy/go-json"
	// "github.com/segmentio/kafka-go"
)

type chatUseCase struct {
	timeout       time.Duration
	chatRepo      r.ChatRepository
	utilU         r.UtilUseCase
	grupoU        r.GrupoUseCase
	salaU         r.SalaUseCase
	wsServer      *ws.WsServer
	conversationU r.ConversationUseCase

	// kafkaW           *kafka.Writer
}

func NewUseCase(timeout time.Duration, charRepo r.ChatRepository, utilU r.UtilUseCase,
	grupoU r.GrupoUseCase, conversationU r.ConversationUseCase,
	salaU r.SalaUseCase, wsServer *ws.WsServer) r.ChatUseCase {
	// w := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9094"),
	// 	Topic:    "notification-message-group",
	// 	Balancer: &kafka.LeastBytes{},
	// }
	return &chatUseCase{
		timeout:       timeout,
		chatRepo:      charRepo,
		wsServer:      wsServer,
		grupoU:        grupoU,
		salaU:         salaU,
		conversationU: conversationU,
		// kafkaW:           w,
		utilU: utilU,
	}
}
func (u *chatUseCase) GetChatUnreadMessages(ctx context.Context, d r.RequestChatUnreadMessages) (res []r.Message, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	switch d.TypeChat {
	case r.TypeChatGrupo:
		res, err = u.grupoU.GetChatUnreadMessage(ctx, d.ChatId, d.LastUpdateChat)
		if err != nil {
			u.utilU.LogError("GetChatUnreadMessages_Grupo", "chat_usecase", err.Error())
			return
		}
		return
	}
	return
}

func (u *chatUseCase) PublishMessage(ctx context.Context, msg r.MessagePublishRequest) (res int, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	switch msg.TypeChat {
	case r.TypeChatGrupo:
		log.Println(msg.Message)
		err = u.grupoU.SaveGrupoMessage(ctx, &msg.Message)
		res = msg.Message.Id
		if err != nil {
			u.utilU.LogError("PublisMessage_SaveGrupoMessage", "chat_usecase", err.Error())
			return
		}
		messagePaylaod, err1 := json.Marshal(msg.Message)
		if err != nil {
			u.utilU.LogError("PublisMessage_Marshal", "chat_usecase", err.Error())
			return 0, err1
		}
		event := r.MessageEvent{
			Type:    "message",
			Payload: string(messagePaylaod),
		}
		payload, err2 := json.Marshal(event)
		if err != nil {
			u.utilU.LogError("PublisMessage_Marshal", "chat_usecase", err.Error())
			return 0, err2
		}
		u.wsServer.Publish(payload, msg.ChatId)
	case r.TypeChatInboxEstablecimiento:
		log.Println("Conversation Message", msg.Message)
		err = u.conversationU.SaveMessage(ctx, &msg.Message)
		res = msg.Message.Id
		if err != nil {
			u.utilU.LogError("PublishMessage_SaveConversationMessage", "chat_usecase", err.Error())
			return
		}
		messagePaylaod, err1 := json.Marshal(msg.Message)
		if err != nil {
			u.utilU.LogError("PublisMessage2_Marshal", "chat_usecase", err.Error())
			return 0, err1
		}
		event := r.MessageEvent{
			Type:    "message",
			Payload: string(messagePaylaod),
		}
		payload, err2 := json.Marshal(event)
		if err != nil {
			u.utilU.LogError("PublisMessage2_Marshal", "chat_usecase", err.Error())
			return 0, err2
		}
		u.wsServer.Publish(payload, msg.ChatId)
	case r.TypeChatSala:
		log.Println("Sala Message", msg.Message)
		err = u.salaU.SaveMessage(ctx, &msg.Message)
		res = msg.Message.Id
		if err != nil {
			u.utilU.LogError("PublishMessage_SaveSalaMessage", "chat_usecase", err.Error())
			return
		}
		messagePaylaod, err1 := json.Marshal(msg.Message)
		if err != nil {
			u.utilU.LogError("PublisMessage3_Marshal", "chat_usecase", err.Error())
			return 0, err1
		}
		event := r.MessageEvent{
			Type:    "message",
			Payload: string(messagePaylaod),
		}
		payload, err2 := json.Marshal(event)
		if err != nil {
			u.utilU.LogError("PublisMessage3_Marshal", "chat_usecase", err.Error())
			return 0, err2
		}
		u.wsServer.Publish(payload, msg.ChatId)
	}
	return
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
