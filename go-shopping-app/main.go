package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-fruit-cart/grpcserver"
	"go-fruit-cart/pkg/routes"

	"github.com/gin-gonic/gin"
)

var (
	port     = "8000"
	grpcPort = ":5050"
	errc     = make(chan error)
	grpcErr  = make(chan error)
	done     = make(chan struct{})
)

func main() {

	router := gin.Default()
	routes.HandleAllRoutes(router)
	//router.POST("/user", controllers.NewUser)
	http.Handle("/", router)
	addrString := fmt.Sprintf("localhost:%v", port)
	//router.Run(addrString)

	reviewGRPCHandler := grpcserver.NewGRPCHandler()
	grpcServer := grpcserver.NewGRPCServer(grpcPort, reviewGRPCHandler)

	go func() {
		errc <- router.Run(addrString)
	}()

	go func() {
		fmt.Printf("fruit seller gRPC Server is running on the port: %v \n", grpcPort)
		grpcErr <- grpcServer.ListenAndServe()
	}()

	select {
	case err := <-errc:
		log.Printf("ListenAndServe error: %v", err)
	case err := <-grpcErr:
		log.Printf("(gRPC Server) ListenAndServe error: %v", err)
	case <-done:
		log.Println("shutting down server ...")
	}
	time.AfterFunc(1*time.Second, func() {
		close(done)
		close(errc)
	})

}
