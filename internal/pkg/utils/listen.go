package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
)

func ListenAndServe(port string) error {
	for {
		ln, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Println(err)
		}

		conn, _ := ln.Accept()
		req := GetRequest(conn)
		req = ParsePort(req)

		if req.Secure { // https
			req.Host = req.Host[:len(req.Host)-4]
			conn.Write([]byte("HTTP/1.0 200 Connection established\r\nProxy-agent: curl/7.79.1\r\n\r\n"))

			path, _ := filepath.Abs("")
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

			msg, err := TlsReadMessage(tlsSrv)
			if err != nil {
				log.Fatal(err)
			}

			tlsConnTo, err := tls.Dial("tcp", req.Host+":443", tlsCfg)
			if err != nil {
				log.Fatal(err)
			}

			err = TlsSendMessage(tlsConnTo, msg)
			if err != nil {
				log.Fatal(err)
			}

			answer, err := TlsReadMessage(tlsConnTo)
			if err != nil {
				log.Fatal(err)
			}

			err = TlsSendMessage(tlsSrv, answer)
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

			ReturnResponse(conn, ProxyRequest(connTo, req.Message))
		}
		conn.Close()
		ln.Close()
	}
}

func DBConnect(Username, DBName, Password, DBHost, DBPort string) (*pgx.ConnPool, error) {
	ConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		Username,
		DBName,
		Password,
		DBHost,
		DBPort)

	ConnConfig, err := pgx.ParseConnectionString(ConnStr)
	if err != nil {
		log.Fatalf("Error config: %s", err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     ConnConfig,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatalf("Error %s during connection to database", err)
	}
	return pool, nil
}
