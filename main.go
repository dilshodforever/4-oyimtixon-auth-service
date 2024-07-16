package main

import (
	"fmt"
	"log"

	"github.com/dilshodforever/4-oyimtixon-auth-service/api"
	"github.com/dilshodforever/4-oyimtixon-auth-service/api/handler"
	pb "github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/auth"
	pbu"github.com/dilshodforever/4-oyimtixon-auth-service/genprotos/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	UserConn, err := grpc.NewClient(fmt.Sprintf("localhost%s", ":8085"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while NEwclient: ", err.Error())
	}
	defer UserConn.Close()
	aus:=pb.NewAuthServiceClient(UserConn)
	us:=pbu.NewUserServiceClient(UserConn)

	h := handler.NewHandler(aus, us)
	r := api.NewGin(h)

	fmt.Println("Server started on port:8081")
	err = r.Run(":8081")
	if err != nil {
		log.Fatal("Error while running server: ", err.Error())
	}
}
