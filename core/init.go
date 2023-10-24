package core

import (
	"database/sql"
	_conversationDeliveryWs "message/core/conversation/delivery/ws"
	_conversationDeliveryHttp "message/core/conversation/delivery/http"
	_conversationR "message/core/conversation/repository/pg"
	_conversationU "message/core/conversation/usecase"

	_chatDeliveryHttp "message/core/chat/delivery/http"
	_chatR "message/core/chat/repository/pg"
	_chatU "message/core/chat/usecase"

	_grupoDeliveryHttp "message/core/grupo/delivery/http"
	_grupoR "message/core/grupo/repository/pg"
	_grupoU "message/core/grupo/usecase"

	_utilU "message/core/util/usecase"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(db *sql.DB){
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,echo.HeaderAccessControlAllowCredentials},
	  }))
	// e.Use(middleware.Logger())
	timeoutContext := time.Duration(5) * time.Second

	utilU := _utilU.NewUseCase()
	//Chat
	chatR := _chatR.NewRepository(db)
	chatU := _chatU.NewUseCase(timeoutContext,chatR,utilU)
	_chatDeliveryHttp.NewHandler(e,chatU)
	//Grupo
	grupoR := _grupoR.NewRepository(db)
	grupoU := _grupoU.NewUseCase(timeoutContext,grupoR,utilU)
	_grupoDeliveryHttp.NewHandler(e,grupoU)

	//Conversation
	conversationR := _conversationR.NewRepository(db)
	conversationU := _conversationU.NewUseCase(timeoutContext,conversationR,utilU)
	_conversationDeliveryHttp.NewHandler(e,conversationU)
	_conversationDeliveryWs.NewWsHandler(e,conversationU)

	conversationAR := _conversationR.NewAdminRepository(db)
	conversationAU := _conversationU.NewAdminUseCase(timeoutContext,conversationAR,utilU)
	_conversationDeliveryHttp.NewAdminHandler(e,conversationAU)

	e.Start("0.0.0.0:9091")
}