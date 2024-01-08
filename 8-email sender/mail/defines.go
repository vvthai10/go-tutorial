package mail

type Request struct {
	Subject     string
	Content     string
	To          []string
	Cc          []string
	Bcc         []string
	AttachFiles []string
	Options     interface{}
}

type Response struct {
	Success bool
	Message string
	Data    string
}
