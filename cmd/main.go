package main

import (
	"bytes"
	"crypto/tls"
	"http-proxy/internal/pkg/utils"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	for {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Println(err)
		}
		conn, _ := ln.Accept()

		req := utils.GetRequest(conn)
		ss := string(req.FullMsg)
		ss += ""

		port := utils.ParsePort(req.Host)
		if req.Secure {
			req.Host = req.Host[:len(req.Host)-4]
			conn.Write([]byte("HTTP/1.0 200 Connection established\r\nProxy-agent: Golang-Proxy\r\n\r\n"))

			path, _ := filepath.Abs("")
			err = exec.Command(path+"/gen_cert.sh", req.Host, strconv.Itoa(rand.Int())).Run()
			if err != nil {
				panic(err)
			}

			cert, err := tls.LoadX509KeyPair(path+"/hck.crt", path+"/cert.key")
			if err != nil {
				panic(err)
			}

			tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
			tlsSrv := tls.Server(conn, tlsCfg)

			msg, err := utils.TlsReadMessage(tlsSrv)
			if err != nil {
				panic(err)
			}

			tlsConnTo, err := tls.Dial("tcp", req.Host+":443", tlsCfg)
			if err != nil {
				panic(err)
			}

			err = utils.TlsSendMessage(tlsConnTo, msg)
			if err != nil {
				panic(err)
			}

			answer, err := utils.TlsReadMessage(tlsConnTo)
			if err != nil {
				panic(err)
			}

			if strings.LastIndex(string(answer), "Transfer-Encoding: chunked") != -1 {
				answer = bytes.Replace(answer, []byte("Transfer-Encoding: chunked"), []byte(""), -1)
			}

			err = utils.TlsSendMessage(tlsSrv, answer)
			if err != nil {
				panic(err)
			}

			tlsConnTo.Close()
			tlsSrv.Close()

		} else {
			if port != "" {
				req.Port = ""
			}
			connTo, err := net.Dial("tcp", req.Host+req.Port)
			if err != nil {
				panic(err)
			}

			answer := utils.ProxyRequest(connTo, req.FullMsg)
			utils.ReturnResponse(conn, answer)
		}

		conn.Close()
		ln.Close()
		ln = nil
		conn = nil
	}
}
