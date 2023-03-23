package interfaces

import "context"

type UserStorage interface {
	Register(ctx context.Context, login, password string) (ok bool, uuid string, err error)
	Login(ctx context.Context, login, password string) (ok bool, uuid string, err error)
}
