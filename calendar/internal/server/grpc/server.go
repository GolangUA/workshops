package grpc

type Validator interface {
	Validate(interface{}) error
}

type Server struct {
	valid Validator
}

func NewServer(valid Validator) *Server {
	return &Server{valid: valid}
}
