package main

import (
	"context"
	pb "go-fruit-cart/newproto"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5050"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure() /* , grpc.WithBlock() */)
	if err != nil {
		log.Fatalf("could not connect:%v", err)
	}
	defer conn.Close()
	c := pb.NewFruitCartManagementServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var userCredentials = []string{"rm@fc.com", "cs@fc.com", "as@fc.com"}

	var productCredentials = []string{"7b0f0075-ec62-4171-9265-4cea81b38c64", "d210f293-e8a1-4fa1-b63a-f5da3bbb3d57", "aa7faf53-6001-4b96-a02a-67708d3a72bb"}

	//implementing get user
	for _, email := range userCredentials {
		r, err := c.GetUser(ctx, &pb.UserCredential{Email: email})
		if err != nil {
			log.Fatalf("error in the request:%v", err)
		}

		log.Printf(`New User
		Firstname:%s,
		Lastname:%s,
		Email:%s
		Cartid:%s`, r.GetFirsname(), r.GetLastname(), r.GetEmail(), r.GetCartid())
	}

	//implementing get product details
	for _, id := range productCredentials {
		r, err := c.GetProduct(ctx, &pb.ProductCredential{Productid: id})
		if err != nil {
			log.Fatalf("error in the request:%v", err)
		}

		log.Printf(`Product
		Name:%s,
		Description:%s,
		Amount:%v
		Image:%s`, r.GetProductname(), r.GetProductdescription(), r.GetProductamount(), r.GetProductimage())
	}

	r, err := c.GetUsers(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("error in the request:%v", err)
	}
	log.Printf(`Users->%v`, r.GetUsers())

}

//implementing get all users
/* */
