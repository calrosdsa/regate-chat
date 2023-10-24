package http

import (
	"context"
	"errors"
	"log"
	ws "message/domain/ws"
	"net/http"
	"strconv"

	r "message/domain/repository"

	"github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
)

type wsChatHandler struct{
	wsServer *ws.WsServer
}

func NewHandler(e *echo.Echo,wsServer *ws.WsServer){
	h := &wsChatHandler{
		wsServer: wsServer,
	}
	e.GET("v1/ws/suscribe/chat/",h.SubscribeHandler)
}

func (h *wsChatHandler) SubscribeHandler(c echo.Context)(err error) {
	id ,err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity,r.ResponseMessage{Message: err.Error()})
	}
	profileId ,err := strconv.Atoi(c.QueryParam("profileId"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity,r.ResponseMessage{Message: err.Error()})
	}
	log.Println(id,"--",profileId)
	err = h.wsServer.Subscribe(c.Request().Context(), c.Response(), c.Request(), id,profileId)
	if errors.Is(err, context.Canceled) {
		log.Println(err)
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
	websocket.CloseStatus(err) == websocket.StatusGoingAway {
	return
	}
	if err != nil {
		log.Println(err)
		h.wsServer.Logf("%v", err)
		return
	}
	return
}

// func (h *wsChatHandler) PublishMessage(c echo.Context)(err error) {
// 	id ,err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusUnprocessableEntity,r.ResponseMessage{Message: err.Error()})
// 	}

// }