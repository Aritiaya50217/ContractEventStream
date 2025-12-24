package main

import (
	"log"
	"os"
	"workflow-service/internal/application/usecase"
	"workflow-service/internal/config"
	"workflow-service/internal/infrastructure/kafka"
	"workflow-service/internal/infrastructure/postgres"
	"workflow-service/internal/interface/consumer"
	httpHandler "workflow-service/internal/interface/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. ")
	}

	// Producer
	kafkaCreateProducer := kafka.NewProducer(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_CREATE_TOPIC"),
	)

	kafkaApproveProducer := kafka.NewProducer(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_APPROVE_TOPIC"),
	)

	// Consumer
	kafkaCreateConsumer := consumer.NewWorkflowConsumer(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_CREATE_TOPIC"),
	)

	kafkaApproveConsumer := consumer.NewWorkflowConsumer(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_APPROVE_TOPIC"),
	)

	kafkaCreateConsumer.Start()
	kafkaApproveConsumer.Start()

	db := config.NewPostgresDB()
	repo := postgres.NewWorkflowRepoPg(db)
	createUsecase := usecase.NewCreateWorkflowUsecase(repo, kafkaCreateProducer)
	approveUsecase := usecase.NewApproveWorkflowUsecase(repo, kafkaApproveProducer)

	handler := httpHandler.NewHandler(createUsecase, approveUsecase)

	r := gin.Default()
	r.POST("/workflows/create", handler.Create)
	r.POST("/workflows/approve/:id", handler.Approve)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
