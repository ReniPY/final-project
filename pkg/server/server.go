package server

import (
	"log"
	"net/http"

	"github.com/ReniPY/final-project/pkg/api"
)

func Run() {
	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	port := ":7540"
	log.Printf("Сервер запущен на порту %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
