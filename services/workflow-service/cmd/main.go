package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"workflow-service/internal/application/usecase"
	"workflow-service/internal/config"
	"workflow-service/internal/infrastructure/kafka"
	"workflow-service/internal/infrastructure/postgres"
	"workflow-service/internal/infrastructure/redis"
	httpHandler "workflow-service/internal/interface/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	kafkago "github.com/segmentio/kafka-go"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. ")
	}

	// Database
	db := config.NewPostgresDB()
	repo := postgres.NewWorkflowRepoPg(db)
	auditRepo := postgres.NewAuditRepoPG(db)
	outboxRepo := postgres.NewOutboxRepoPG(db)

	// redis cache
	redisClient := redis.NewRedisClient()
	workflowCache := redis.NewWorkflowCache(redisClient)

	// Kafka Producer & Consumer
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	// Producer เดียวสำหรับ Create + Approve
	producer := kafka.NewProducer(kafkaBroker, topic)

	// Reader สำหรับ consumer
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
		GroupID: "workflow-audit-consumer",
	})

	consumer := kafka.NewWorkflowConsumer(reader, auditRepo)
	// Start consumer in goroutine
	go consumer.Start()

	outboxWorker := kafka.NewOutboxWorker(outboxRepo, *producer, topic)
	go outboxWorker.Start()

	createUsecase := usecase.NewCreateWorkflowUsecase(repo, producer, workflowCache)
	approveUsecase := usecase.NewApproveWorkflowUsecase(repo, workflowCache, outboxRepo)
	getWorkflowUsecase := usecase.NewGetWorkflowUsecase(repo, workflowCache)

	handler := httpHandler.NewHandler(createUsecase, approveUsecase, getWorkflowUsecase)

	r := gin.Default()
	r.POST("/workflows/create", handler.Create)
	r.POST("/workflows/approve/:id", handler.Approve)
	r.GET("/workflows/:id", handler.GetWorkflow)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Workflow service running on port", port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}
}
