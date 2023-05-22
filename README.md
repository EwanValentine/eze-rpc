# Eze RPC

Is a simple RPC library for Golang, which uses binary encoding/decoding based on  a DSL for type safety. 

Eze uses code generation to generate the boilerplate code.

This project is purely for a "fun" because I'm sick of gRPC, and am imagining a world where there's something easier to use. So... don't use it in production. If you do, then it's your own silly fault if something goes horribly wrong.

What's funny is, it basically doesn't solve any problem that gRPC has, other than being smaller. But it's smaller because it barely functions. But it killed a few hours, and who knows, maybe there's enough other people out there who CBA with gRPC anymore, that this could take off and actual adults could take over the development of this.


## Installation

```bash
$ go install github.com/EwanValentine/eze-rpc
```

## Usage

### Define an Eze schema

```
package users

service UserService {
	GetUser(request: GetUserRequest): User
	CreateUser(request: CreateUserRequest): User
}

struct CreateUserRequest {
	Name: String
}

struct User {
	ID: String
	Name: String
}

struct GetUserRequest {
	ID: String
}
```

### Run the code generate command
```bash
eze generate --input=example.eze --output=.
```

### Implement the generated interface
```go
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
```

### Register the server

```go
srv := users.NewServer()
srv.RegisterService("UserService", &UserService{})
if err := srv.Serve(":8080"); err != nil {
    panic(err)
}
```

### Call the server

```go
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
```
