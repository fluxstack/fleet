package usecase

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"github.com/weflux/fastapi/module/user/config"
	"github.com/weflux/fastapi/module/user/entity"
	"github.com/weflux/fastapi/module/user/repo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewAuth(users repo.Users, cfg config.Config) *Auth {
	return &Auth{users: users, cfg: cfg}
}

type Auth struct {
	users repo.Users
	cfg   config.Config
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

func (uc *Auth) Login(ctx context.Context, req *LoginRequest) (*LoginResult, error) {
	if req.Login == "" || req.Password == "" {
		return nil, errors.New("请输入账号和密码")
	}

	user, err := uc.users.GetByUsername(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if user != nil && user.ID > 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return nil, errors.New("用户名或密码错误")
		}
		return uc.makeTokens(user.ID)
	}
	user, err = uc.users.GetByPhone(ctx, req.Login)
	if user != nil && user.ID > 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return nil, errors.New("手机号或密码错误")
		}
		return uc.makeTokens(user.ID)
	}

	return nil, errors.New("账号或密码错误")
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (uc *Auth) makeTokens(id entity.UID) (*LoginResult, error) {
	var accessToken, refreshToken string
	var err error
	now := time.Now()
	accessExpired := now.Add(24 * time.Hour)
	refreshExpired := now.Add(30 * 24 * time.Hour)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: int64(id),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "weflux",
			Subject:   cast.ToString(id),
			ExpiresAt: jwt.NewNumericDate(accessExpired),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: int64(id),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "weflux",
			Subject:   cast.ToString(id),
			ExpiresAt: jwt.NewNumericDate(refreshExpired),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	accessToken, err = at.SignedString([]byte(uc.cfg.Auth.JWTSecret))
	if err != nil {
		return nil, err
	}
	refreshToken, err = rt.SignedString([]byte(uc.cfg.Auth.JWTSecret))
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  24 * 60 * 60,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: 30 * 24 * 60 * 60,
	}, nil
}
