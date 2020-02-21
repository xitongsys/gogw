package schema

type RegisterRequest struct {
	SourceAddr string
	ToPort int
	Direction string
	Protocol string
	Description string
}

type RegisterResponse struct {
	ClientId string
	Status string
}