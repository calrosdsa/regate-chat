package http

import (
	// "log"

	"log"
	r "message/domain/repository"
	"net/http"
	"strconv"

	// "strconv"

	"github.com/labstack/echo/v4"
	_jwt "message/domain/util"
)

type ConversationHandler struct {
	conversationUcase r.ConversationUseCase
}

func NewHandler(e *echo.Echo, conversationUcase r.ConversationUseCase) {
	handler := ConversationHandler{
		conversationUcase: conversationUcase,
	}
	e.GET("v1/conversation/messages/:id/",handler.GetMessages)
	// e.GET("v1/conversation/messages/:id/",handler.GetConversationMessages)
	e.GET("v1/conversations/",handler.GetConversations)
	e.POST("v1/conversation/sync-messages/",handler.SyncMessage)
	e.GET("v1/conversation/get-id/:establecimientoId/",handler.GetOrCreateConversation)
}
func (h *ConversationHandler)GetOrCreateConversation(c echo.Context)(err error){
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	claims, err := _jwt.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	establecimientoId,err := strconv.Atoi(c.Param("establecimientoId"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	conversationId,err := h.conversationUcase.GetOrCreateConversation(ctx,establecimientoId,claims.ProfileId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct{
		Id int `json:"id"`
	}{
		Id: conversationId,
	}
	return c.JSON(http.StatusOK,response)
}

func (h *ConversationHandler)SyncMessage(c echo.Context)(err error){
	log.Println("SYNC")
	var data []r.Inbox
	err = c.Bind(&data)
	if err != nil {
		log.Println("SYNC",err)
		return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	for _,message := range data{
		log.Println("SYNC MESSAGE",message)
		err := h.conversationUcase.SaveMessage(ctx,&message)
		if err != nil {
		    log.Println("SYNC error",err)
			// return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
		}
		// log.Println("SYNC",err)
	}
	return c.JSON(http.StatusOK, data)
}


func (h *ConversationHandler) GetMessages(c echo.Context) (err error){
	// auth := c.Request().Header["Authorization"][0]
	// token := _jwt.GetToken(auth)
	// claims, err := _jwt.ExtractClaims(token)
	// if err != nil {
		// log.Println(err)
		// return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	// }
	log.Println("GETTING MESSAGES")
	page,err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
	}
	log.Println("PAGE",page)
	id ,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,nextPage,err := h.conversationUcase.GetMessages(ctx,id,int16(page),20)
	if err != nil {
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct {
		Page int16 `json:"nextPage"`
		Results  []r.Inbox `json:"results"`
	}{
		Page: nextPage,
		Results: res,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *ConversationHandler) GetConversations(c echo.Context) (err error){
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	claims, err := _jwt.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := h.conversationUcase.GetConversations(ctx,claims.ProfileId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// func (h *ConversationHandler)GetConversationMessages(c echo.Context)(err error){
// 	page,err := strconv.Atoi(c.QueryParam("page"))
// 	// if err != nil {
// 	// 	log.Println(err)

// 	// 	return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
// 	// }
// 	log.Println("PAGE",page)
// 	id ,err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusConflict, r.ResponseMessage{Message: err.Error()})
// 	}
// 	ctx := c.Request().Context()
// 	res,err := h.conversationUcase.GetConversationMessages(ctx,id,page)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
// 	}
// 	return c.JSON(http.StatusOK, res)
// }
