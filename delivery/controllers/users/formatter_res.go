package user

import "altastore/entities"

type RegisterUserResponseFormat struct {
	Message string          `json:"message"`
	Data    []entities.User `json:"data"`
}

type LoginUserResponseFormat struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type GetUsersResponseFormat struct {
	Message string          `json:"message"`
	Data    []entities.User `json:"data"`
}

type GetUserResponseFormat struct {
	Message string        `json:"message"`
	Data    entities.User `json:"data"`
}

type PutUserResponseFormat struct {
	Message string        `json:"message"`
	Data    entities.User `json:"data"`
}

type DeleteUserResponseFormat struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
