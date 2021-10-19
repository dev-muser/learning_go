package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dev-muser/learning_go/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	hello_handler := handlers.NewHello(l)

	servemux := http.NewServeMux() // HTTP request multiplexer
	servemux.Handle("/hello", hello_handler)
	fmt.Println("Server is up and running.")
	http.ListenAndServe(":7777", servemux)

}
