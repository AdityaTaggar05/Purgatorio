package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidToken = errors.New("token has expired or has been already used")
	ErrUserNotFound = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidPasswordFormat = errors.New("invalid password format")
)

