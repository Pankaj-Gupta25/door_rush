package api

import (
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/sachinggsingh/PDTS-Go/pb/parcel"
	"github.com/sachinggsingh/PDTS/parcel-service/config"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/api/grpcapi"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/api/restapi"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/messaging"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/repository"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/service"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
	"github.com/sachinggsingh/PDTS/parcel-service/pkg/db"
)

type Server struct {
	DB     *db.Database
	Logger *utils.Logger
	ENV    *config.ENV
}

func NewServer(env *config.ENV, logger *utils.Logger, DB *db.Database) *Server {
	return &Server{
		ENV:    env,
		Logger: logger,
		DB:     DB,
	}
}

func (s *Server) RunServer() {
	r := gin.New()

	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Add a simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok from parcel-service"})
	})

	parcelRepo := repository.NewParcelRepository(s.DB, s.Logger)
	parcelStore := repository.NewParcelStore()

	// Kafka Producer
	var kafkaProducer messaging.KafkaProducer
	if s.ENV.KAFKA_BROKERS != "" {
		brokers := []string{s.ENV.KAFKA_BROKERS}
		kafkaProducer = messaging.NewKafkaProducer(brokers, "parcel.events")
		s.Logger.Info("Kafka producer initialized for topic: parcel.events")
	}

	parcelService := service.NewParcelService(parcelRepo, parcelStore, kafkaProducer, s.Logger)

	parcelHandler := restapi.NewParcelHandler(parcelService, s.Logger, r)

	parcelHandler.SetupRoutes(s.ENV.JWT_SECRET)

	// Start gRPC server
	go func() {
		grpcPort := "50051"
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			s.Logger.Error("Failed to listen for gRPC: " + err.Error())
			return
		}
		grpcServer := grpc.NewServer()
		handler := grpcapi.NewParcelGrpcHandler(parcelService, s.Logger)
		pb.RegisterParcelServiceServer(grpcServer, handler)

		s.Logger.Info("gRPC Server is running on port " + grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Error("Failed to serve gRPC: " + err.Error())
		}
	}()

	s.Logger.Info("Server is running on port " + s.ENV.PORT)
	if err := r.Run(":" + s.ENV.PORT); err != nil {
		s.Logger.Error("Failed to start server: " + err.Error())
	}
}
