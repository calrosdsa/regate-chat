package ws

import (
	// "encoding/json"
	// "log"
	"net/http"
	// "strconv"

	// "net/http"
	r "message/domain/repository"

	// _jwt "message/domain/util"

	"github.com/labstack/echo/v4"
)

type SalaHandler struct {
	grupoU r.SalaUseCase
}

func NewHandler(e *echo.Echo,salaUseCase r.SalaUseCase){
	// go H.Run(salaUseCase)
	handler := SalaHandler{
		grupoU: salaUseCase,
	}
	// e.GET("v1/ws/chat-grupo",handler.ChatGrupo)
	e.GET("v1/chat/grupo/unread-messages/",handler.GetUnreadMessages)
	// e.GET("v1/chat/grupo/unread-messages/:chatId/",handler.GetChatUnreadMessages)
}



func (h *SalaHandler)GetUnreadMessages(c echo.Context)(err error){
	// auth := c.Request().Header["Authorization"][0]
	// token := _jwt.GetToken(auth)
	// claims, err := _jwt.ExtractClaims(token)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	// }
	// page,err := strconv.Atoi(c.QueryParam("page"))
	// if err != nil {
	// 	page = 1
	// }
	// size := 20
	// ctx := c.Request().Context()
	// res,nextPage,err := h.grupoU.GetUnreadMessages(ctx,claims.ProfileId,int16(page),int8(size))
	// if err != nil {
	// 	log.Println("SYNC error",err)
	// 	return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	// }
	// response := struct {
	// 	Page int16 `json:"page"`
	// 	Results  []r.MessageGrupo `json:"results"`
	// }{
	// 	Page: nextPage,
	// 	Results: res,
	// }
	return c.JSON(http.StatusOK, "response")
}

// func (ws *SalaHandler) ChatGrupo(c echo.Context) (err error) {
// 	casoId := c.QueryParam("id")
// 	log.Println("chat grupo",casoId)
// 	ServeWs(c.Response(), c.Request(), casoId)
// 	return nil
// }

