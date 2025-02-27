package api

import (
	"context"
	"github.com/weflux/fastapi/errors"
	"github.com/weflux/fastapi/module/user/usecase"
)

func NewAPI(uc *usecase.Auth) *API {
	return &API{uc: uc}
}

type API struct {
	uc *usecase.Auth
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginReply struct {
	errors.APIError
	Data *usecase.LoginResult `json:"data"`
}

func (api *API) Login(ctx context.Context, req *usecase.LoginRequest) (*LoginReply, error) {
	res, err := api.uc.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return &LoginReply{
		Data: res,
	}, nil
}
