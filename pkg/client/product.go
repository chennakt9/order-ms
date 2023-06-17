package client

import (
	"context"
	"log"

	"github.com/chennakt9/order-ms/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func InitProductServiceClient(url string) ProductServiceClient {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// fmt.Println("conn str is", conn)

	if err != nil {
		log.Fatalln("Could not connect:", err)
	}

	c := ProductServiceClient{
		Client: pb.NewProductServiceClient(conn),
	}

	return c
}

func (c *ProductServiceClient) FindOne(productId int64) (*pb.FindOneResponse, error) {
	req := &pb.FindOneRequest{
		Id: productId,
	}

	return c.Client.FindOne(context.Background(), req)
}

func (c *ProductServiceClient) DecreaseStock(productId int64, quantity int64, orderId int64) (*pb.DecreaseStockResponse, error) {
	req := &pb.DecreaseStockRequest {
		Id: productId,
		OrderId: orderId,
		Quantity: quantity,
	}

	return c.Client.DecreaseStock(context.Background(), req)
}