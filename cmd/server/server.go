package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/declantraynor/trunc/storage"
	"github.com/declantraynor/trunc/web"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	ExternalHost string `split_words:"true"`
	Port         int
	RedisHost    string `split_words:"true"`
	RedisPort    int    `split_words:"true"`
}

func main() {
	var config Configuration
	if err := envconfig.Process("trunc", &config); err != nil {
		log.Fatal(err.Error())
	}

	urlBuilder, err := web.NewRandomURLBuilder(config.ExternalHost)
	if err != nil {
		log.Fatalf("Check EXTERNAL_HOST config value: %v", err)
	}

	store, err := storage.NewRedisStore(config.RedisHost, config.RedisPort)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer store.Disconnect()

	service := &web.Service{
		URLBuilder: urlBuilder,
		Store:      store,
	}

	router := mux.NewRouter()
	router.HandleFunc("/shorten", service.Shorten).Methods("POST")
	router.PathPrefix("/").HandlerFunc(service.Redirect).Methods("GET")
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", config.Port),
			router,
		),
	)
}
