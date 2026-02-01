package models

type SuccessResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type RequestDelete struct {
	ID int `json:"id"`
}
