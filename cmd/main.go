package main

import (
	"fmt"
	"log"
	"net"

	"github.com/chennakt9/order-ms/pkg/client"
	"github.com/chennakt9/order-ms/pkg/config"
	"github.com/chennakt9/order-ms/pkg/db"
	"github.com/chennakt9/order-ms/pkg/pb"
	"github.com/chennakt9/order-ms/pkg/service"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	fmt.Println(c)

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	productSvc := client.InitProductServiceClient(c.ProductServiceUrl)

	// fmt.Println("Order svc on", c.Port)

	s := service.Server{
		H: h,
		ProductSvc: productSvc,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to service", err)
	}
}