package main

import (
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	pDelivery "http-proxy/internal/pkg/proxy/delivery"
	pRepo "http-proxy/internal/pkg/proxy/repo"
	"http-proxy/internal/pkg/scaner/delivery"
	"http-proxy/internal/pkg/scaner/repo"
	"http-proxy/internal/pkg/scaner/usecase"
	"http-proxy/internal/pkg/utils"
	"log"
	"net/http"
)

type tomlConfig struct {
	Title string
	Web   webServer   `toml:"web-server"`
	Proxy proxyServer `toml:"proxy-server"`
	DB    database    `toml:"database"`
}

type webServer struct {
	Host string `toml:"host"`
	Port string
}

type proxyServer struct {
	Host string
	Port string
}

type database struct {
	DbName   string
	Username string
	Password string
	Host     string
	Port     string
}

func main() {
	var conf tomlConfig
	if _, err := toml.DecodeFile("./conf.toml", &conf); err != nil {
		log.Fatal(err)
	}

	db, err := utils.DBConnect(conf.DB.Username, conf.DB.DbName, conf.DB.Password, conf.DB.Host, conf.DB.Port)
	if err != nil {
		log.Fatal(err)
	}

	newRepo := pRepo.NewRepoPostgres(db)
	proxyServ := pDelivery.NewProxyServer(newRepo, ":"+conf.Proxy.Port)
	go func() {
		log.Fatal(proxyServ.ListenAndServe())
	}()

	muxRoute := mux.NewRouter()

	rRepo := repo.NewRepoPostgres(db)
	rUsecase := usecase.NewRepoUsecase(rRepo)
	handler := delivery.NewRepeaterHandler(rUsecase)

	repeater := muxRoute.PathPrefix("/api/v1").Subrouter()
	{
		repeater.HandleFunc("/requests", handler.AllRequests).Methods(http.MethodGet)
		repeater.HandleFunc("/requests/{id}", handler.GetRequest).Methods(http.MethodGet)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":"+conf.Web.Port, muxRoute))
}
