package main

import (
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"http-proxy/internal/pkg/proxy/delivery"
	"http-proxy/internal/pkg/proxy/repo"
	"http-proxy/internal/pkg/proxy/usecase"
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

	go func() {
		log.Fatal(utils.ListenAndServe(conf.Proxy.Port))
	}()

	muxRoute := mux.NewRouter()

	db, err := utils.DBConnect(conf.DB.Username, conf.DB.DbName, conf.DB.Password, conf.DB.Host, conf.DB.Port)
	if err != nil {
		log.Fatal(err)
	}

	fRepo := repo.NewRepoPostgres(db)
	fUsecase := usecase.NewRepoUsecase(fRepo)
	fHandler := delivery.NewForumHandler(fUsecase)

	forum := muxRoute.PathPrefix("/api/v1").Subrouter()
	{
		forum.HandleFunc("/forum/create", fHandler.AllRequest).Methods(http.MethodGet)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":"+conf.Web.Port, muxRoute))
}
