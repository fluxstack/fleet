package entity

import "time"

type UID int64

type User struct {
	// 主键
	ID UID `json:"id,omitempty"`
	// 用户名
	Name string `json:"name,omitempty"`
	// 昵称
	DisplayName string `json:"display_name,omitempty"`
	// 头像
	AvatarURL string `json:"avatar_url,omitempty"`
	// 手机号
	Phone string `json:"phone,omitempty"`
	// 邮箱
	Email string `json:"email,omitempty"`
	// 密码hash
	PasswordHash string `json:"-"`
	// 创建时间戳
	CreatedAt time.Time `json:"created_at,omitempty"`
	// 更新时间戳
	LastUpdatedAt time.Time `json:"last_updated_at,omitempty"`
	// 是否删除，1-是，0-否
	IsDeleted int8 `json:"is_deleted,omitempty"`
	// 权限标记，0-普通用户，1-普通管理员，2-超管
	Role int32 `json:"role,omitempty"`
	// 是否锁定，1-是，0-否
	IsLocked int8 `json:"is_locked,omitempty"`
}
