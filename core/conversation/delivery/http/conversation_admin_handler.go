package http

import (
	// "log"

	"log"
	r "message/domain/repository"
	"net/http"
	// "strconv"

	// "strconv"

	"github.com/labstack/echo/v4"
	_jwt "message/domain/util"
)

type ConversationAdminHandler struct {
	conversationAdminUseCase r.ConversationAdminUseCase
}

func NewAdminHandler(e *echo.Echo, conversationAdminUseCase r.ConversationAdminUseCase) {
	handler := ConversationAdminHandler{
		conversationAdminUseCase: conversationAdminUseCase,
	}
	// e.GET("v1/conversation/messages/:id/",handler.GetMessages)
	// e.GET("v1/conversation/messages/:id/",handler.GetConversationMessages)
	// e.GET("v1/conversations/",handler.GetConversations)
	e.GET("v1/conversations/:uuid/",handler.GetConversationsEstablecimiento)
}

func (h *ConversationAdminHandler)GetConversationsEstablecimiento(c echo.Context)(err error){
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	_, err = _jwt.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	uuid := c.Param("uuid")
	// if err != nil {
		// log.Println("SYNC",err)
		// return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
	// }
	ctx := c.Request().Context()
	res,err := h.conversationAdminUseCase.GetConversationsEstablecimiento(ctx,uuid)
	if err != nil {
		log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
