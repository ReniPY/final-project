package server

import (
	"log"
	"net/http"

	"github.com/ReniPY/final-project/pkg/api"
)

func Run() {
	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
