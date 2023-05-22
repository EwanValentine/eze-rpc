package users

import (
    "net"
    "fmt"
    "encoding/gob"
	"reflect"
)

type UserService interface {
	GetUser(request *GetUserRequest) (*User, error)
	CreateUser(request *CreateUserRequest) (*User, error)
}


type CreateUserRequest struct {
	Name string
}
type User struct {
	ID string
	Name string
}
type GetUserRequest struct {
	ID string
}

func NewUserServiceClient(conn func() (net.Conn, error)) *UserServiceClient {
	return &UserServiceClient{
		conn: conn,
	}
}

type UserServiceClient struct {
	conn func() (net.Conn, error)
}

func NewConnection(addr string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

type Request struct {
    ServiceName string
    MethodName  string
    Arg         interface{}
}
func (c *UserServiceClient) GetUser(request *GetUserRequest) (*User, error) {
	conn, err := c.conn()
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	gob.Register(&GetUserRequest{})
	gob.Register(&User{})

	r := &Request{
		ServiceName: "UserService",
		MethodName:  "GetUser",
		Arg:         request,
	}
	err = encoder.Encode(r)
	if err != nil {
		return nil, fmt.Errorf("error encoding arg: %v", err)
	}

	var response *User
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response type: %v", err)
	}

	return response, nil
}
func (c *UserServiceClient) CreateUser(request *CreateUserRequest) (*User, error) {
	conn, err := c.conn()
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	gob.Register(&CreateUserRequest{})
	gob.Register(&User{})

	r := &Request{
		ServiceName: "UserService",
		MethodName:  "CreateUser",
		Arg:         request,
	}
	err = encoder.Encode(r)
	if err != nil {
		return nil, fmt.Errorf("error encoding arg: %v", err)
	}

	var response *User
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response type: %v", err)
	}

	return response, nil
}

type Server struct {
    services map[string]interface{}
}

func NewServer() *Server {
    return &Server{
        services: make(map[string]interface{}),
    }
}

func (s *Server) RegisterService(name string, service interface{}) {
    s.services[name] = service
}

func (s *Server) Serve(address string) error {
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

    var request Request
    err := decoder.Decode(&request)
    if err != nil {
        fmt.Println("Error decoding request:", err)
        return
    }

    service, ok := s.services[request.ServiceName]
    if !ok {
        fmt.Println("No such service:", request.ServiceName)
        return
    }

    results := reflect.ValueOf(service).MethodByName(request.MethodName).Call([]reflect.Value{
        reflect.ValueOf(request.Arg),
    })

    if len(results) != 2 || !results[1].IsNil() {
        fmt.Println("Error calling method:", request.MethodName)
        return
    }

    response := results[0].Interface()
    err = encoder.Encode(response)
    if err != nil {
        fmt.Println("Error encoding response:", err)
        return
    }
}
