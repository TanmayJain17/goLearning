package grpcserver

import (
	"context"
	pb "go-fruit-cart/newproto"
	"go-fruit-cart/pkg/models"
	"log"
)

// GRPCHandler implements the rating management gRPC interface
type FruitCartManagementServer struct {
	pb.UnimplementedFruitCartManagementServiceServer
}

// NewGRPCHandler returns new rating management service gRPC Handler
func NewGRPCHandler() *FruitCartManagementServer {
	return &FruitCartManagementServer{}
}

func (s *FruitCartManagementServer) GetUser(ctx context.Context, in *pb.UserCredential) (*pb.UserDetails, error) {

	log.Printf("Received user %v", in.GetEmail())
	email, firstname, lastname, cartid, err := models.GrpcFindTheUser(in.GetEmail())
	if err != nil {
		log.Fatalf("UserDetails not found error occured-> %v", err.Error())
	}
	return &pb.UserDetails{Firsname: firstname, Lastname: lastname, Email: email, Cartid: cartid}, nil
}

func (s *FruitCartManagementServer) GetProduct(ctx context.Context, in *pb.ProductCredential) (*pb.ProductDetails, error) {

	log.Printf("Received user %v", in.GetProductid())
	name, description, price, image, err := models.GrpcFindProduct(in.GetProductid())
	if err != nil {
		log.Fatalf("ProductDetails not found error occured-> %v", err.Error())
	}
	return &pb.ProductDetails{Productname: name, Productdescription: description, Productamount: price, Productimage: image}, nil
}

func (s *FruitCartManagementServer) GetUsers(ctx context.Context, in *pb.Empty) (*pb.AllUsers, error) {
	array := []models.Userdetails{}
	array, err := models.GrpcGetAllUsers()
	if err != nil {
		log.Fatalf("ProductDetails not found error occured-> %v", err.Error())
	}
	var usersArray []*pb.UserDetails
	for _, user := range array {
		var theUser = pb.UserDetails{}
		theUser.Firsname = user.Firstname
		theUser.Lastname = user.Lastname
		theUser.Email = user.Email
		theUser.Cartid = user.Cartid
		usersArray = append(usersArray, &theUser)
	}
	return &pb.AllUsers{Users: usersArray}, nil
}
