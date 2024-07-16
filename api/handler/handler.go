package handler

import (
	pb "github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/auth"
	pbu "github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/user"
)

type Handler struct {
	Auth pb.AuthServiceClient
	User pbu.UserServiceClient
}

func NewHandler(auth pb.AuthServiceClient, user pbu.UserServiceClient) *Handler {
	return &Handler{
		Auth: auth,
		User: user,
	}

}
