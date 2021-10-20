package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/dev-muser/learning_go/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	products_handler := handlers.NewProducts(l)
	servemux := mux.NewRouter()
	// Request comes to server
	// It gets picked up by the router
	// See that is a PUT / POST request (because subrouter.Use(middleware))
	// send request to subrouter and execute the middleware. If it passes, then
	// goes to the handle func.
	getRouter := servemux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", products_handler.GetProducts)

	putRouter := servemux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", products_handler.UpdateProducts)
	putRouter.Use(products_handler.MiddlewareProductValidation) //middleware validation applied

	postRouter := servemux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", products_handler.AddProduct)
	postRouter.Use(products_handler.MiddlewareProductValidation) //middleware validation applied

	fmt.Println("Server is up and running.")
	server := &http.Server{
		Addr:         ":7777",
		Handler:      servemux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
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
	//trap sigterm or interupt
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt) //broadcast the message on the channel
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan //Block until there's a message available to be consumed
	l.Println("Received terminate, graceful shutdown", sig)

	timeout_context, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout_context)

}
