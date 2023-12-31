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

type salaUcase struct {
	timeout  time.Duration
	salaRepo r.SalaRepository
	kafkaW   *kafka.Writer
	utilU    r.UtilUseCase
}

func NewUseCase(timeout time.Duration, salaRepo r.SalaRepository, utilU r.UtilUseCase) r.SalaUseCase {
	w := &kafka.Writer{
		Addr:     kafka.TCP(viper.GetString("kafka.host")),
		Topic:    "notification-message-sala",
		Balancer: &kafka.LeastBytes{},
	}
	return &salaUcase{
		timeout:  timeout,
		salaRepo: salaRepo,
		kafkaW:   w,
		utilU:    utilU,
	}
}
func (u *salaUcase) GetChatUnreadMessages(ctx context.Context, chatId int, lastUpdate string) (res []r.Message, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err = u.salaRepo.GetChatUnreadMessages(ctx, chatId, lastUpdate)
	return
}

func (u *salaUcase)DeleteMessage(ctx context.Context,id int)(err error){
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	err = u.salaRepo.DeleteMessage(ctx,id)
	return
}

func (u *salaUcase) GetUsers(ctx context.Context,d r.RequestUsersGroupOrRoom) (res []r.UsersGroupOrRoom, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res,err = u.salaRepo.GetUsers(ctx, d)
	return
}

func (u *salaUcase) SaveMessage(ctx context.Context, d *r.Message) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	err = u.salaRepo.SaveMessage(ctx, d)
	if err != nil {
		u.utilU.LogError("SaveMessage", "grupo_usecase", err.Error())
		return
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
