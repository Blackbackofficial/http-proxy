package utils

import (
	"bytes"
	"fmt"
	"http-proxy/internal/models"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseHost(headers []byte) string {
	i := strings.LastIndex(string(headers), "Host:")
	i = i + len("Host:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	host := string(headers)[i : i+j]
	return host
}

func ParseSecure(headers []byte) bool {
	i := strings.LastIndex(string(headers), "CONNECT")
	if i == -1 {
		return false
	}
	return true
}

func ParsePort(req models.Request) models.Request {
	re := regexp.MustCompile(`(?m):\d+`)
	res := re.FindAllString(req.Host, -1)
	if res == nil {
		return req
	}
	req.Port = ""
	return req
}

func ParseLength(headers []byte) int {
	fmt.Print("Received headers:\n", string(headers))
	i := strings.LastIndex(string(headers), "Content-Length:")
	i = i + len("Content-Length:") + 1
	j := strings.Index(string(headers)[i:], "\r")
	l, _ := strconv.Atoi(string(headers)[i : i+j])
	return l
}

func GetRequest(conn net.Conn) models.Request {
	var bMessage []byte
	var bBody []byte

	for {
		bArr := make([]byte, 10, 10)
		n, err := conn.Read(bArr)
		if err != nil || n == 0 {
			return models.Request{}
		}
		bMessage = append(bMessage, bArr...)
		if strings.Contains(string(bMessage), "\r\n\r\n") {
			i := strings.LastIndex(string(bMessage), "\r\n\r\n") + len("\r\n\r\n")
			bBody = append(bBody, bMessage[i:len(bMessage)]...)
			bMessage = bMessage[:i]
			break
		}
	}

	if strings.LastIndex(string(bMessage), "Proxy-Connection: Keep-Alive\r\n") != -1 {
		bMessage = bytes.Replace(bMessage, []byte("Proxy-Connection: Keep-Alive\r\n"), []byte(""), -1)
	}

	if strings.LastIndex(string(bMessage), "http://") != -1 {
		i := strings.LastIndex(string(bMessage), "http://")
		j := i + len("http://")
		for string(bMessage)[j] != '/' {
			j++
		}
		bMessage = bytes.Replace(bMessage, bMessage[i:j], []byte(""), 1)
	}

	l := ParseLength(bMessage)
	fmt.Print("Start receiving body\n")
	for {
		if l == 0 {
			break
		}
		bArr := make([]byte, l, l)
		n, err := conn.Read(bArr)
		if err != nil || n == 0 {
			return models.Request{}
		}
		bBody = append(bBody, bArr...)
		if n > l-len(bBody) {
			break
		}
	}
	fmt.Print("Received body:\n", string(bBody))

	fullMsg := make([]byte, 0, 10)
	fullMsg = append(fullMsg, bMessage...)
	fullMsg = append(fullMsg, bBody...)
	return models.Request{Port: ":80", Message: fullMsg, Secure: ParseSecure(bMessage), Host: ParseHost(bMessage)}
}

func ReturnResponse(conn net.Conn, answer []byte) {
	sentBytes := 0
	for sentBytes < len(answer) {
		n, err := conn.Write(answer)
		if err != nil {
			fmt.Println("Error in sending data", n)
		}
		sentBytes += n
	}
}

func ProxyRequest(conn net.Conn, msg []byte) []byte {
	sentBytes := 0
	for sentBytes < len(msg) {
		n, err := conn.Write(msg)
		if err != nil {
			fmt.Println("Error in sending data", n)
		}
		sentBytes += n
	}
	answer := GetRequest(conn)
	return answer.Message
}

func TlsReadMessage(conn net.Conn) ([]byte, error) {
	var msg []byte
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	for {
		bArr := make([]byte, 1024, 1024)
		n, err := conn.Read(bArr)
		if n == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		msg = append(msg, bArr...)
	}
	msg = append(msg, []byte("\r\n\r\n")...)
	return msg, nil
}

func TlsSendMessage(conn net.Conn, msg []byte) error {
	bytesSent := 0
	for bytesSent < len(msg) {
		n, err := conn.Write(msg)
		if err != nil {
			return err
		}
		bytesSent += n
	}
	return nil
}