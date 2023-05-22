package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EwanValentine/eze-rpc/example/users"
)

type UserService struct{}

func (s *UserService) GetUser(user *users.GetUserRequest) (*users.User, error) {
	return &users.User{
		ID:   user.ID,
		Name: "Ewan",
	}, nil
}

func (s *UserService) CreateUser(user *users.CreateUserRequest) (*users.User, error) {
	return &users.User{
		ID:   "123",
		Name: user.Name,
	}, nil
}

func main() {
	go func() {
		srv := users.NewServer()
		srv.RegisterService("UserService", &UserService{})
		if err := srv.Serve(":8080"); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second * 1)

	getConnection := users.NewConnection(":8080")

	client := users.NewUserServiceClient(getConnection)
	response, err := client.GetUser(&users.GetUserRequest{
		ID: "123",
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(response.Name)

	createResponse, err := client.CreateUser(&users.CreateUserRequest{
		Name: "Ewan",
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createResponse.ID)
}
