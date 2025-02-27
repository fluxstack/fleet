package repoimpl

import (
	"context"
	"database/sql"
	"errors"
	"github.com/weflux/fastapi/module/user/entity"
	"github.com/weflux/fastapi/module/user/repo"
	"github.com/weflux/fastapi/module/user/repoimpl/model"
	"github.com/weflux/fastapi/storage/database"
	"time"
)

type users struct {
	*database.Client
}

func newUser(v *model.User) *entity.User {
	return &entity.User{
		ID:            entity.UID(v.ID),
		Name:          v.Name,
		DisplayName:   v.DisplayName,
		AvatarURL:     v.AvatarURL,
		Phone:         v.Phone,
		Email:         v.Email,
		PasswordHash:  v.PasswordHash,
		CreatedAt:     time.UnixMilli(v.CreatedAt),
		LastUpdatedAt: time.UnixMilli(v.LastUpdatedAt),
		IsDeleted:     v.IsDeleted,
		Role:          v.Role,
		IsLocked:      v.IsLocked,
	}
}
func (impl *users) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	v := &model.User{}
	if err := impl.DB.NewSelect().Model(v).Where("name = ?", username).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	return newUser(v), nil
}

func (impl *users) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	v := &model.User{}
	if err := impl.DB.NewSelect().Model(v).Where("phone = ?", phone).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	return newUser(v), nil
}

func (impl *users) Create(ctx context.Context, user *repo.UserCreate) error {
	_, err := impl.DB.NewInsert().Model(&model.User{
		ID:            int64(user.ID),
		Name:          user.Name,
		DisplayName:   user.DisplayName,
		Phone:         user.Phone,
		PasswordHash:  user.PasswordHash,
		CreatedAt:     user.CreatedAt.UnixMilli(),
		LastUpdatedAt: user.CreatedAt.UnixMilli(),
	}).Exec(ctx)

	return err
}

func NewUsers(
	dbClients *database.Client,
) repo.Users {
	return &users{
		dbClients,
	}
}
