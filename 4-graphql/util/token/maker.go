package token

type Maker interface {
	CreateToken(arg PayloadParams) (string, error)
	VerifyToken(token string) (*Payload, error)
}