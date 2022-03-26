package models

type Request struct {
	FullMsg []byte
	Secure  bool
	Host    string
	Port    string
}

type Response struct {
	FullMsg []byte
}
