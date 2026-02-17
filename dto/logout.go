package dto

type LogoutReq struct {
	Token string `json:"token"`
}

type LogoutRes struct {
	Message string `json:"message"`
}
