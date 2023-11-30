package repository

type ResponseMessage struct {
	Message string `json:"message"`
}

type UsersGroupOrRoom struct {
	Id              int      `json:"id"`
	TypeChat        TypeChat `json:"type_chat"`
	IsAdmin         bool     `json:"is_admin,omitempty"`
	ProfileId       int      `json:"profile_id"`
	ProfileName     string   `json:"nombre"`
	ProfileApellido *string  `json:"apellido,omitempty"`
	ProfilePhoto    *string  `json:"profile_photo,omitempty"`
	IsOut           bool     `json:"is_out"`
	ParentId        int      `json:"parent_id"`
}

type TypeEntity int8

const (
	ENTITY_NONE            = 0
	ENTITY_SALA            = 1
	ENTITY_GRUPO           = 2
	ENTITY_ACCOUNT         = 3
	ENTITY_BILLING         = 4
	ENTITY_RESERVA         = 5
	ENTITY_ESTABLECIMIENTO = 6
	ENTITY_URI             = 7
	ENTITY_SALA_COMPLETE   = 8
	ENTITY_INVITATION      = 9
)
