package usecase

import (
	"context"
	"github.com/goccy/go-json"
	"log"
	r "message/domain/repository"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type grupoUcase struct {
	timeout   time.Duration
	grupoRepo r.GrupoRepository
	kafkaW    *kafka.Writer
	utilU     r.UtilUseCase
}

func NewUseCase(timeout time.Duration, grupoRepo r.GrupoRepository, utilU r.UtilUseCase) r.GrupoUseCase {
	w := &kafka.Writer{
		Addr:     kafka.TCP(viper.GetString("kafka.url")),
		Topic:    "notification-message-group",
		Balancer: &kafka.LeastBytes{},
	}
	return &grupoUcase{
		timeout:   timeout,
		grupoRepo: grupoRepo,
		kafkaW:    w,
		utilU:     utilU,
	}
}
func (u *grupoUcase) GetChatUnreadMessage(ctx context.Context, chatId int64, lastUpdate string) (res []r.Message, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	res, err = u.grupoRepo.GetChatUnreadMessage(ctx, chatId, lastUpdate)
	if err != nil {
		u.utilU.LogError("GetChatUnreadMessage", "grupo_usecase", err.Error())
		return
	}
	return
}

func (u *grupoUcase) GetUnreadMessages(ctx context.Context, profileId int, page int16,
	size int8) (res []r.Message, nextPage int16, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	page = u.utilU.PaginationValues(page)
	res, err = u.grupoRepo.GetUnreadMessages(ctx, profileId, page, int8(size))
	if err != nil {
		u.utilU.LogError("SaveGrupoMessage", "grupo_usecase", err.Error())
	}
	nextPage = u.utilU.GetNextPage(int8(len(res)), int8(size), page+1)
	go u.UpdateUserGrupoLastTimeUpdateMessage(context.Background(), profileId)
	return
}

func (u *grupoUcase) UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	err = u.grupoRepo.UpdateUserGrupoLastTimeUpdateMessage(ctx, profileId)
	if err != nil {
		u.utilU.LogError("UpdateUserGrupoLastTimeUpdateMessage", "grupo_ucase", err.Error())
	}
	return
}

func (u *grupoUcase) SaveGrupoMessage(ctx context.Context, d *r.Message) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer func() {
		cancel()
	}()
	err = u.grupoRepo.SaveGrupoMessage(ctx, d)
	if err != nil {
		u.utilU.LogError("SaveGrupoMessage", "grupo_usecase", err.Error())
		return
	}
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
