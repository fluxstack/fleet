package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int64  `bun:"id,pk"`
	Name         string `bun:"name"`
	DisplayName  string `bun:"display_name"`
	AvatarURL    string `bun:"avatar_url"`
	Phone        string `bun:"phone"`
	Email        string `bun:"email"`
	PasswordHash string `bun:"password_hash"`
	CreatedAt    int64  `bun:"created_at"`
	//CreatedBy     int64  `bun:"created_by"`
	LastUpdatedAt int64 `bun:"last_updated_at"`
	//LastUpdatedBy int64  `bun:"last_updated_by"`
	Role      int32 `bun:"role"`
	IsDeleted int8  `bun:"is_deleted"`
	IsLocked  int8  `bun:"is_locked"`
}
