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

type ChatHandler struct {
	chatUseCase r.ChatUseCase
}

func NewHandler(e *echo.Echo, chatUseCase r.ChatUseCase) {
	handler := ChatHandler{
		chatUseCase: chatUseCase,
	}
	// e.GET("v1/conversation/messages/:id/",handler.GetMessages)
	// e.GET("v1/conversation/messages/:id/",handler.GetConversationMessages)
	// e.GET("v1/conversations/",handler.GetConversations)
	e.GET("v1/chats/",handler.GetChatsUser)
	e.POST("v1/chat/publish/message/",handler.PublishMessage)
	e.POST("v1/chat/unread-messages/",handler.GetChatUnreadMessages)
	e.POST("v1/chat/delete/message/",handler.DeleteMessage)
	e.GET("v1/chat/deleted/messages/:id/",handler.GetDeletedMessages)
	e.POST("v1/chat/users/",handler.GetUsers)

}

func (h *ChatHandler)GetUsers(c echo.Context)(err error){
	// auth := c.Request().Header["Authorization"][0]
	// token := _jwt.GetToken(auth)
	// _, err = _jwt.ExtractClaims(token)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	// }
	var data r.RequestUsersGroupOrRoom
	err = c.Bind(&data)
	log.Println("GET USERS",data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := h.chatUseCase.GetUsers(ctx,data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	log.Println(res,"users")
	return c.JSON(http.StatusOK,res)
}


func (h *ChatHandler)DeleteMessage(c echo.Context)(err error){
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	_, err = _jwt.ExtractClaims(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	var data r.DeleteMessageRequet
	err = c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	err = h.chatUseCase.DeleteMessage(ctx,data)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,nil)
}

func (h *ChatHandler) GetChatUnreadMessages(c echo.Context)(err error) {
	var data r.RequestChatUnreadMessages
	err = c.Bind(&data)
	if err != nil {
		log.Println("SYNC error",err)
		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
	}
	log.Println(data)
	ctx := c.Request().Context()
	res,err := h.chatUseCase.GetChatUnreadMessages(ctx,data)
	if err != nil {
		// log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK,res)
}

func (h *ChatHandler) GetDeletedMessages(c echo.Context)(err error) {
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
	}
	ctx := c.Request().Context()
	res,err := h.chatUseCase.GetDeletedMessages(ctx,id)
	if err != nil {
		// log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct {
		Ids []int `json:"ids"`
	}{
		Ids: res,
	}
	return c.JSON(http.StatusOK,response)
}


// func (h *ChatHandler) SyncMessages(c echo.Context)(err error) {
// 	var data []r.MessagePublishRequest
// 	err = c.Bind(&data)
// 	if err != nil {
// 		log.Println("SYNC error",err)
// 		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
// 	}
// 	log.Println(data)
// 	ctx := c.Request().Context()
// 	for i :=0;i < len(data);i++{
// 		h.chatUseCase.PublishMessage(ctx,data)
// 	}
// 	// if err != nil {
// 		// log.Println("SYNC error",err)
// 		// return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
// 	// }
// 	return nil
// }

// func (h *ChatHandler) SharedMessage(c echo.Context)(err error) {
// 	var data []r.MessagePublishRequest
// 	err = c.Bind(&data)
// 	if err != nil {
// 		log.Println("SYNC error",err)
// 		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
// 	}
// 	log.Println(data)
// 	ctx := c.Request().Context()
// 	res,err := h.chatUseCase.PublishMessage(ctx,data)
// 	if err != nil {
// 		// log.Println("SYNC error",err)
// 		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
// 	}
// 	response := struct {
// 		Id int `json:"id"`
// 	}{
// 		Id: res,
// 	}
// 	return c.JSON(http.StatusOK,response)
// }

func (h *ChatHandler) PublishMessage(c echo.Context)(err error) {
	var data r.MessagePublishRequest
	err = c.Bind(&data)
	if err != nil {
		log.Println("SYNC error",err)
		return c.JSON(http.StatusUnprocessableEntity, r.ResponseMessage{Message: err.Error()})
	}
	log.Println(data)
	ctx := c.Request().Context()
	res,err := h.chatUseCase.PublishMessage(ctx,data)
	if err != nil {
		// log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct {
		Id int `json:"id"`
	}{
		Id: res,
	}
	return c.JSON(http.StatusOK,response)
}

func (h *ChatHandler)GetChatsUser(c echo.Context)(err error){
	log.Println("GETTING CHATS")
	auth := c.Request().Header["Authorization"][0]
	token := _jwt.GetToken(auth)
	claims, err := _jwt.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, r.ResponseMessage{Message: err.Error()})
	}
	page,err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	size := 20
	ctx := c.Request().Context()
	res,nextPage,err := h.chatUseCase.GetChatsUser(ctx,claims.ProfileId,int16(page),int8(size))
	if err != nil {
		log.Println("SYNC error",err)
		return c.JSON(http.StatusBadRequest, r.ResponseMessage{Message: err.Error()})
	}
	response := struct {
		Page int16 `json:"page"`
		Results  []r.Chat `json:"results"`
	}{
		Page: nextPage,
		Results: res,
	}
	return c.JSON(http.StatusOK, response)
}
