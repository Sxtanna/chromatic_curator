package backend

import "context"

type Backend interface {
	RoleBackend
}

type RoleBackend interface {
	GetRole(ctx context.Context, user string) (string, error)

	SetRole(ctx context.Context, user string, role string) error
}
