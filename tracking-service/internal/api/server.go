package api

import (
	"context"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/sachinggsingh/PDTS-Go/pb/tracking"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/config"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/api/grpcapi"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/api/restapi"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/messaging"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/repository"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/service"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/websockets"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/pkg/db"
)

type Server struct {
	Router *gin.Engine
	logger *utils.Logger
	db     *db.Database
}

func NewServer() *Server {
	database := db.NewDatabase()
	if err := database.ConnectToDatabase(); err != nil {
		panic(err)
	}

	return &Server{
		Router: gin.Default(),
		logger: utils.Log,
		db:     database,
	}
}

func (s *Server) StartServer() {
	r := s.Router

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK from tracking-service",
		})
	})

	hub := websockets.NewHub()
	go hub.Run()

	repo := repository.NewTrackingRepository(s.db, s.logger)
	svc := service.NewTrackingService(repo, hub, s.logger)
	handler := restapi.NewTrackingHandler(svc, hub, s.logger)

	// Kafka Integration
	env := config.LoadENV()
	if env.KAFKA_BROKERS != "" {
		brokers := []string{env.KAFKA_BROKERS} // Assuming single broker for now or comma-separated
		consumer := messaging.NewKafkaConsumer(brokers, "parcel.events", "tracking-group", svc, s.logger)
		go consumer.Start(context.Background())
		s.logger.Info("Kafka consumer initialized for topic: parcel.events")
	}

	handler.SetupRoutes(r, env.JWT_SECRET)

	// Start gRPC server
	go func() {
		grpcPort := "50052"
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			s.logger.Error("Failed to listen for gRPC: " + err.Error())
			return
		}
		grpcServer := grpc.NewServer()
		handler := grpcapi.NewTrackingGrpcHandler(svc, s.logger)
		pb.RegisterTrackingServiceServer(grpcServer, handler)

		s.logger.Info("gRPC Server is running on port " + grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error("Failed to serve gRPC: " + err.Error())
		}
	}()

	s.logger.Info("Server started on port " + env.PORT)
	r.Run(":" + env.PORT)
}
