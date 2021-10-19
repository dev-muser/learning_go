package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	//implementing graceful shutdown
	// for usecases like big file upload or db transaction, 
	// could have the risk of disconecting my client, 
	// not allowing to finish the work.

	go func() {
		err := server.ListenAndServe() // will going to block so put inside a goroutine
		if err != nil {
			l.Fatal(err)
		}
	}()

	//but because is still imediately start to shutdown
	// going to use os.signal to register for certain signal notification

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt) //broadcast the message on the channel
	signal.Notify(sigChan, os.Kill)
	
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	
	timeout_context,  _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout_context)

}
