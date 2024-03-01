package main

import (
	validation "av-send-email/api/config"
	av_repository "av-send-email/api/pkg/av/repository"
	av_service "av-send-email/api/pkg/av/service"
	"av-send-email/api/routes"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

const (
	AWS_REGION = "us-east-1"
	AWS_BUCKET = "avlogssvg"
)
func main(){
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	AWS_KEY := os.Getenv("AWS_KEY")
	AWS_SECRET_KEY := os.Getenv("AWS_SECREY_KEY")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(AWS_REGION),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AWS_KEY, AWS_SECRET_KEY, "")),
		config.WithClientLogMode(aws.LogRetries | aws.LogRequest | aws.LogResponse),
	)
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	sesClient := sesv2.NewFromConfig(cfg)
	avRepo := av_repository.NewRepository(s3Client, sesClient)
	myValidator := validation.NewValidator()
	avService := av_service.NewService(avRepo)

    app := fiber.New(fiber.Config{
		ServerHeader:  "Fiber scan av",
		AppName: "avscan-email",
        // Global custom error handler
        ErrorHandler: func(c fiber.Ctx, err error) error {
            return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
                Success: false,
                Message: err.Error(),
            })
        },
    })

	app.Use(cors.New())
	api := app.Group("/api")
	routes.Routes(api, avService, myValidator)
	  // Catch-all route
	  app.Use(func(c fiber.Ctx) error {
        // Send a custom error message
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Route not found",
        })
    })
	log.Fatal(app.Listen(":8088"))
}
