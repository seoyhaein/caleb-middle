package apiV1

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/samples/ch02/productinfo/go/server/ecommerce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server is used to implement ecommerce/product_info.
// etcd grpc_proxy.go, maintenace.go, rpc.pb.go 에서 regisger 하는 부분 차이점을 살펴보자.
type server struct {
	productMap map[string]*ecommerce.Product
}

func NewProductServer() ecommerce.ProductInfoServer {
	return &server{
		productMap: make(map[string]*ecommerce.Product),
	}
}

// AddProduct implements ecommerce.AddProduct
func (s *server) AddProduct(ctx context.Context,
	in *ecommerce.Product) (*ecommerce.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*ecommerce.Product)
	}
	s.productMap[in.Id] = in
	log.Printf("Product %v : %v - Added.", in.Id, in.Name)
	return &ecommerce.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct implements ecommerce.GetProduct
func (s *server) GetProduct(ctx context.Context, in *ecommerce.ProductID) (*ecommerce.Product, error) {
	product, exists := s.productMap[in.Value]
	if exists && product != nil {
		log.Printf("Product %v : %v - Retrieved.", product.Id, product.Name)
		return product, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}
