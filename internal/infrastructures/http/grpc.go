package http

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/controllers"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"

	pkg "5g-v2x-data-management-service/pkg/api"
	"log"
	"net"

	"go.uber.org/dig"
	"google.golang.org/grpc"
)

var grpcServe *grpc.Server

type GRPCServer struct {
	GRPCGateway GRPCServerGateway
	config      *config.Config
}

type GRPCServerGateway struct {
	dig.In
	ControllerGateway controllers.ControllerGateway
}

func NewGRPCServer(gateway GRPCServerGateway, config *config.Config) *GRPCServer {
	return &GRPCServer{
		GRPCGateway: gateway,
		config:      config,
	}
}

func init() {
	lr := logrus.New()
	lr.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logger := logrus.NewEntry(lr)
	grpcLogrus.ReplaceGrpcLogger(logger)

	grpcServe = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpcLogrus.StreamServerInterceptor(logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpcLogrus.UnaryServerInterceptor(logger),
		)),
	)

}

func (g *GRPCServer) registerServer() {
	pkg.RegisterAccidentServiceServer(grpcServe, g.GRPCGateway.ControllerGateway.AccidentController)
	pkg.RegisterDrowsinessServiceServer(grpcServe, g.GRPCGateway.ControllerGateway.DrowsinessController)
	pkg.RegisterCarServiceServer(grpcServe, g.GRPCGateway.ControllerGateway.CarController)
}

func (g *GRPCServer) Start() error {
	g.registerServer()
	lis, err := net.Listen("tcp", g.config.ServiceAddress)
	if err != nil {
		return err
	}
	log.Println("[DATA-MANAGEMENT-SERVICE] Listening ", g.config.ServiceAddress)
	return grpcServe.Serve(lis)
}
