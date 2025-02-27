package repo

import (
	"context"
	"github.com/weflux/fastapi/module/user/entity"
	"time"
)

type Users interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Create(ctx context.Context, user *UserCreate) error
}

type UserCreate struct {
	ID           entity.UID `json:"id"`
	Phone        string     `json:"phone"`
	PasswordHash string     `json:"password_hash"`
	Name         string     `json:"name"`
	DisplayName  string     `json:"display_name"`
	CreatedAt    time.Time  `json:"created_at"`
}
