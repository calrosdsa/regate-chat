package http

import (
	// "log"

	"log"
	r "message/domain/repository"
	"net/http"
	"strconv"

	// "strconv"

	// "strconv"

	_jwt "message/domain/util"

	"github.com/labstack/echo/v4"
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
	e.GET("v1/conversations/messages-count/:uuid/",handler.GetConversationsMessagesCount)
	e.GET("v1/conversation/messages/:id/",handler.GetMessages)
}

func (h *ConversationAdminHandler) GetMessages(c echo.Context) (err error) {
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	_, err = _jwt.ExtractClaimsAdmin(token)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	log.Println("PAGE MESSAGE",page)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	size := 20
	ctx := c.Request().Context()
	res, nextPage, err := h.conversationAdminUseCase.GetMessages(ctx, id, int16(page), int8(size))
	if err != nil {
		log.Println("SYNC error", err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct {
		Page    int16       `json:"page"`
		Results []r.MessageWithReply `json:"results"`
	}{
		Page:    nextPage,
		Results: res,
	}
	return c.JSON(http.StatusOK, response)
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

func (h *ConversationAdminHandler)GetConversationsMessagesCount(c echo.Context)(err error){
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	_, err = _jwt.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	uuid := c.Param("uuid")
	ctx := c.Request().Context()
	res,err := h.conversationAdminUseCase.GetConversationsMessagesCount(ctx,uuid)
	if err != nil {
		log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
