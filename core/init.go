package core

import (
	"database/sql"
	"log"
	_conversationDeliveryHttp "message/core/conversation/delivery/http"
	_conversationR "message/core/conversation/repository/pg"
	_conversationU "message/core/conversation/usecase"

	_chatDeliveryHttp "message/core/chat/delivery/http"
	_chatR "message/core/chat/repository/pg"
	_chatU "message/core/chat/usecase"

	_grupoDeliveryHttp "message/core/grupo/delivery/http"
	_grupoR "message/core/grupo/repository/pg"
	_grupoU "message/core/grupo/usecase"

	_salaDeliveryHttp "message/core/sala/delivery/http"
	_salaR "message/core/sala/repository/pg"
	_salaU "message/core/sala/usecase"


	_wsChatDeliveryHttp "message/core/wschat/delivery/http"
	// _wsChatU "message/core/wschat/usecase"

	_utilU "message/core/util/usecase"
	ws "message/domain/ws"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func Init(db *sql.DB){

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,echo.HeaderAccessControlAllowCredentials},
	  }))
	// e.Use(middleware.Logger())
	timeoutContext := time.Duration(15) * time.Second

	//Ws 
	chatWsServer := &ws.WsServer{
		SubscriberMessageBuffer: 16,
		Logf:                    log.Printf,
		Subscribers:             make(map[int]map[int]*ws.Subscriber),
		PublishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	
	utilU := _utilU.NewUseCase()
	//WsChat
	// wsChatU := _wsChatU.NewUseCase(chatWsServer,utilU)
	_wsChatDeliveryHttp.NewHandler(e,chatWsServer)

	
	//Grupo
	grupoR := _grupoR.NewRepository(db)
	grupoU := _grupoU.NewUseCase(timeoutContext,grupoR,utilU)
	_grupoDeliveryHttp.NewHandler(e,grupoU)
	//Sala
	salaR := _salaR.NewRepository(db)
	salaU := _salaU.NewUseCase(timeoutContext,salaR,utilU)
	_salaDeliveryHttp.NewHandler(e,salaU)

	//Conversation
	conversationR := _conversationR.NewRepository(db)
	conversationU := _conversationU.NewUseCase(timeoutContext,conversationR,utilU)
	_conversationDeliveryHttp.NewHandler(e,conversationU)
	
	conversationAR := _conversationR.NewAdminRepository(db)
	conversationAU := _conversationU.NewAdminUseCase(timeoutContext,conversationAR,utilU)
	_conversationDeliveryHttp.NewAdminHandler(e,conversationAU)
	
	//Chat
	chatR := _chatR.NewRepository(db)
	chatU := _chatU.NewUseCase(timeoutContext,chatR,utilU,grupoU,conversationU,salaU,chatWsServer)
	_chatDeliveryHttp.NewHandler(e,chatU)

	e.Start("0.0.0.0:9091")
}