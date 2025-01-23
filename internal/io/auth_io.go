package io

type AuthInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthOutput struct {
	Created bool   `json:"created"`
	Token   string `json:"token"`
}
