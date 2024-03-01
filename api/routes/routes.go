package routes

import (
	validation "av-send-email/api/config"
	"av-send-email/api/handlers"
	av "av-send-email/api/pkg/av/service"

	"github.com/gofiber/fiber/v3"
)

func Routes(app fiber.Router, service av.Service, validator *validation.XValidator){
	app.Post("/scanAv", handlers.ScanAv(service, validator))
}