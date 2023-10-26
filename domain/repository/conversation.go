package repository

import "context"

type EstablecimientoConversation struct {
	Name           string  `json:"name"`
	Apellido       *string `json:"apellido"`
	Photo          *string `json:"photo"`
	ConversationId int     `json:"conversation_id"`
	ProfileId      int     `json:"profile_id"`
	ParentId   int `json:"parent_id"`
}

type ConversationAdminRepository interface {
	GetConversationsEstablecimiento(ctx context.Context, uuid string) ([]EstablecimientoConversation, error)
}
type ConversationAdminUseCase interface {
	GetConversationsEstablecimiento(ctx context.Context, uuid string) ([]EstablecimientoConversation, error)
}

type ConversationRepository interface {
	SaveMessage(ctx context.Context, d *Message) (err error)
	GetMessages(ctx context.Context, id int, page int16, size int8) ([]Inbox, error)
	GetConversations(ctx context.Context, id int) ([]Conversation, error)

	GetOrCreateConversation(ctx context.Context, id int, profileId int) (conversationId int, err error)
}

type ConversationUseCase interface {
	SaveMessage(ctx context.Context, d *Message) (err error)
	GetMessages(ctx context.Context, id int, page int16, size int8) (res []Inbox, nextPage int16, err error)
	GetConversations(ctx context.Context, id int) ([]Conversation, error)

	GetOrCreateConversation(ctx context.Context, id int, profileId int) (conversationId int, err error)
}

type Conversation struct {
	Id                   int    `json:"id"`
	EstablecimientoName  string `json:"establecimiento_name"`
	EstablecimientoId    int    `json:"establecimiento_id"`
	EstablecimientoPhoto string `json:"establecimiento_photo"`
}

type Inbox struct {
	Id             int        `json:"id"`
	SenderId       int        `json:"sender_id"`
	ConversationId int        `json:"conversation_id"`
	Content        string     `json:"content"`
	CreatedAt      string     `json:"created_at,omitempty"`
	ReplyTo        *int       `json:"reply_to"`
	Reply          ReplyInbox `json:"reply"`
}

type ReplyInbox struct {
	Id             int    `json:"id"`
	SenderId       int    `json:"sender_id"`
	ConversationId int    `json:"conversation_id"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at,omitempty"`
	ReplyTo        *int   `json:"reply_to"`
}

type ConversationEvent struct {
	Type    string `json:"type"`
	Message Inbox  `json:"message"`
	// Sala    SalaData     `json:"sala,omitempty"`
}
