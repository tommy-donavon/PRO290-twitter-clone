package data

import "github.com/go-playground/validator"

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l *Login) Validate() error {
	validator := validator.New()
	return validator.Struct(l)
}
