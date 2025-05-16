package backend

import "context"

type Backend interface {
	RoleBackend
}

type RoleBackend interface {
	GetRole(ctx context.Context, guild string, user string) (string, error)

	SetRole(ctx context.Context, guild string, user string, role string) error
}
