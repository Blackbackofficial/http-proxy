package models

type Request struct {
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Scheme  string            `json:"scheme"`
	Host    string            `json:"host"`
	Path    string            `json:"path"`
	Cookies string            `json:"cookies"`
	Header  map[string]string `json:"header,omitempty"`
	Body    string            `json:"body,omitempty"`
}

type Response struct {
	Id        int               `json:"id"`
	RequestId int               `json:"request_id"`
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Cookies   string            `json:"cookies"`
	Header    map[string]string `json:"header"`
	Body      string            `json:"body"`
}

type Cookies struct {
	Key   string
	Value string
}

type Error struct {
	Message string `json:"message"`
}
