package usecase

import (
	// "encoding/json"
	// "context"
	r "message/domain/repository"
	ws "message/domain/ws"
	"time"
) 

type wsChatUseCase struct {
	wsServer *ws.WsServer
	utilU r.UtilUseCase
	timeout time.Duration
}

func NewUseCase(wsServer *ws.WsServer,utilU r.UtilUseCase,timeout time.Duration) ws.WsChatUseCase{
	return &wsChatUseCase{
		wsServer: wsServer,
		utilU: utilU,
		timeout: timeout,
	}
}

// func (u *wsChatUseCase) PublishMessage(ctx int,msg ws.MessagePublishRequest){
// 	ctx,cancel := context.WithTimeout(ctx,u.timeout)
// 	defer cancel()
// 	switch msg.TypeChat {
// 	case r.TypeChatGrupo:
// 		u.wsServer.Publish(msg)
		
// 	}
// }