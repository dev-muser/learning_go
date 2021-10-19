package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dev-muser/learning_go/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	hello_handler := handlers.NewHello(l)
	goodbye_handler := handlers.NewGoodbye(l)

	servemux := http.NewServeMux() // HTTP request multiplexer
	servemux.Handle("/hello", hello_handler)
	servemux.Handle("/goodbye", goodbye_handler)

	fmt.Println("Server is up and running.")
	server := &http.Server{
		Addr: ":7777",
		Handler: servemux,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	server.ListenAndServe()

}
