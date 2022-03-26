package models

type Request struct {
	Message []byte
	Secure  bool
	Host    string
	Port    string
}

type Response struct {
	Msg []byte
}
