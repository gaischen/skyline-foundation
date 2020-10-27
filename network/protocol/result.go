package protocol

//common result
type Result interface {
	Success() bool
	Message() string
	Code() int
	Result() interface{}
}



