package service

import (
	"context"
	"net/http"

	"github.com/chennakt9/order-ms/pkg/client"
	"github.com/chennakt9/order-ms/pkg/db"
	"github.com/chennakt9/order-ms/pkg/models"
	"github.com/chennakt9/order-ms/pkg/pb"
)

type Server struct {
	pb.OrderServiceServer
	H db.Handler
	ProductSvc client.ProductServiceClient
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.OrderSvcNoParam) (*pb.OrderSvcHealthCheckResponse, error) {
	return &pb.OrderSvcHealthCheckResponse{
		Message: "Order service is up",
	}, nil
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error){
	product, err := s.ProductSvc.FindOne(req.ProductId)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if product.Status >= http.StatusNotFound {
		return &pb.CreateOrderResponse{Status: product.Status, Error: product.Error}, nil
	} else if product.Data.Stock < req.Quantity {
		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: "Stock too less"}, nil
	}

	order := models.Order {
		Price: product.Data.Price,
		ProductId: product.Data.Id,
		UserId: req.UserId,
	}

	s.H.DB.Create(&order)

	res, err := s.ProductSvc.DecreaseStock(req.ProductId, req.Quantity, order.Id)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if res.Status == http.StatusConflict {
		s.H.DB.Delete(&models.Order{}, order.Id)

		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

  return &pb.CreateOrderResponse{
    Status: http.StatusCreated,
    Id: order.Id,
  }, nil
}

