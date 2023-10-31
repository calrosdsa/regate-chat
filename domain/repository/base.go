package repository

type ResponseMessage struct {
	Message string `json:"message"`
}

type UsersGroupOrRoom struct {
	Id              int     `json:"id"`
	IsAdmin         bool    `json:"is_admin"`
	ProfileId       int     `json:"profile_id"`
	ProfileName     string  `json:"nombre"`
	ProfileApellido *string `json:"apellido"`
	ProfilePhoto    *string `json:"profile_photo"`
	IsOut           bool    `json:"is_out"`
}
