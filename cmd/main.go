package main

import (
	"crypto/tls"
	"fmt"
	"http-proxy/internal/pkg/utils"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
)

func main() {
	for {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Println(err)
		}

		conn, _ := ln.Accept()
		req := utils.GetRequest(conn)
		req = utils.ParsePort(req)

		if req.Secure { // https
			req.Host = req.Host[:len(req.Host)-4]
			conn.Write([]byte("HTTP/1.0 200 Connection established\r\nProxy-agent: curl/7.79.1\r\n\r\n"))

			path, _ := filepath.Abs("")
			fmt.Print(path, "HIIIII")
			err = exec.Command(path+"/certs/gen_cert.sh", req.Host, strconv.Itoa(rand.Int())).Run()
			if err != nil {
				log.Fatal(err)
			}

			cert, err := tls.LoadX509KeyPair(path+"/certs/nck.crt", path+"/certs/cert.key")
			if err != nil {
				log.Fatal(err)
			}

			tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
			tlsSrv := tls.Server(conn, tlsCfg)

			msg, err := utils.TlsReadMessage(tlsSrv)
			if err != nil {
				log.Fatal(err)
			}

			tlsConnTo, err := tls.Dial("tcp", req.Host+":443", tlsCfg)
			if err != nil {
				log.Fatal(err)
			}

			err = utils.TlsSendMessage(tlsConnTo, msg)
			if err != nil {
				log.Fatal(err)
			}

			answer, err := utils.TlsReadMessage(tlsConnTo)
			if err != nil {
				log.Fatal(err)
			}

			err = utils.TlsSendMessage(tlsSrv, answer)
			if err != nil {
				log.Fatal(err)
			}

			tlsConnTo.Close()
			tlsSrv.Close()
		} else { // http
			connTo, err := net.Dial("tcp", req.Host+req.Port)
			if err != nil {
				log.Fatal(err)
			}

			utils.ReturnResponse(conn, utils.ProxyRequest(connTo, req.Message))
		}

		conn.Close()
		ln.Close()
	}
}
