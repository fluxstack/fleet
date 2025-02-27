package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/weflux/fastapi/errors"
)

type API struct {
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginReply struct {
	errors.APIError
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (api *API) Login(ctx context.Context, req LoginRequest) (*LoginReply, error) {
	return &LoginReply{
		Data: struct {
			Token string `json:"token"`
		}{
			Token: uuid.NewString(),
		},
	}, nil
}
