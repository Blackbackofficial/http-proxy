package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func HeadToStr(headers http.Header) string {
	var stringHeaders string
	for key, values := range headers {
		for _, value := range values {
			stringHeaders += key + " " + value + "\n"
		}
	}
	return stringHeaders
}

func StrToHeader(headers string) map[string]string {
	h := make(map[string]string)
	for _, header := range strings.Split(headers, "\n") {
		if len(header) < 2 {
			continue
		}
		str := strings.SplitN(header, " ", 2)
		h[str[0]] = str[1]
	}
	return h
}

func GenTLSConf(host, URL string) (tls.Config, error) {
	path, _ := filepath.Abs("")
	err := exec.Command(path+"/certs/gen_cert.sh", host, strconv.Itoa(rand.Int())).Run()
	if err != nil {
		log.Fatal(err)
		return tls.Config{}, err
	}

	tlsCert, err := tls.LoadX509KeyPair(path+"/certs/nck.crt", path+"/certs/cert.key")
	if err != nil {
		log.Println("error loading pair", err)
		return tls.Config{}, err
	}

	return tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ServerName:   URL,
	}, nil
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
