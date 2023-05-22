package generator

import (
	"bytes"
	"fmt"
	"text/template"
)

const codeTemplate = `package {{.Package}}

import (
    "net"
    "fmt"
    "encoding/gob"
	"reflect"
)

type {{.Service.Name}} interface {
	{{- range .Service.Methods}}
	{{.Name}}({{.Arg}} *{{.ArgType}}) (*{{.Ret}}, error)
	{{- end}}
}

{{range .Structs}}
type {{.Name}} struct {
	{{- range $field, $type := .Fields}}
	{{$field}} {{mapType $type}}
	{{- end}}
}
{{- end}}

func New{{.Service.Name}}Client(conn func() (net.Conn, error)) *{{.Service.Name}}Client {
	return &{{.Service.Name}}Client{
		conn: conn,
	}
}

type {{.Service.Name}}Client struct {
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

{{- range .Service.Methods}}
func (c *{{$.Service.Name}}Client) {{.Name}}({{.Arg}} *{{.ArgType}}) (*{{.Ret}}, error) {
	conn, err := c.conn()
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %v", err)
	}

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	gob.Register(&{{.ArgType}}{})
	gob.Register(&{{.Ret}}{})

	r := &Request{
		ServiceName: "{{$.Service.Name}}",
		MethodName:  "{{.Name}}",
		Arg:         {{.Arg}},
	}
	err = encoder.Encode(r)
	if err != nil {
		return nil, fmt.Errorf("error encoding arg: %v", err)
	}

	var response *{{.Ret}}
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response type: %v", err)
	}

	return response, nil
}
{{- end}}

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
`

// This is because I opted to use capitalised primitive types in the DSL for no other reason that they look
// slightly nicer to me, so now I have to pointlessly map between the two. If I was being serious about this for
// any reason, other than me being unemployed and having nothing better to do, I probably wouldn't do this.
func mapType(dslType string) string {
	switch dslType {
	case "String":
		return "string"
	case "Int":
		return "int"
	case "Float":
		return "float64"
	case "Bool":
		return "bool"
	default:
		return dslType
	}
}

func GenerateCode(dsl DSL) string {
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"mapType": mapType,
	}
	tmpl, err := template.New("code").Funcs(funcMap).Parse(codeTemplate)
	if err != nil {
		fmt.Println("Template parsing error:", err)
		return ""
	}
	err = tmpl.Execute(&buf, dsl)
	if err != nil {
		fmt.Println("Template execution error:", err)
		return ""
	}
	return buf.String()
}
