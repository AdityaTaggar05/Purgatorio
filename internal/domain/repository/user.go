package repository

import "context"

type UserRepository interface {
	CreateUser(context.Context, string, string) error
}
