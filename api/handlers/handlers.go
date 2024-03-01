package handlers

import (
	validation "av-send-email/api/config"
	av "av-send-email/api/pkg/av/service"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)
func ScanAv(s av.Service, validator *validation.XValidator) fiber.Handler{
	
	return func(c fiber.Ctx) error{
		c.Accepts("application/json")
		body := c.Body()
		 p:= new(validation.ScanAvRequest)
		 err := json.Unmarshal(body, p)
		 if err != nil {
            return errors.New("invalid Request Body")
        }
		fmt.Printf("Parsed JSON: %+v\n", p.Paths)

        if errs := validator.Validate(p); len(errs) > 0 && errs[0].Error {
            errMsgs := make([]string, 0)

            for _, err := range errs {
                errMsgs = append(errMsgs, fmt.Sprintf(
                    "[%s]: '%v' | Needs to implement '%s'",
                    err.FailedField,
                    err.Value,
                    err.Tag,
                ))
            }

            return &fiber.Error{
                Code:    fiber.ErrBadRequest.Code,
                Message: strings.Join(errMsgs, " and "),
            }
        }

		result, errAv := s.ClamAvScan(p.Paths)
		if errAv != nil{
			return errAv
		}

		fmt.Printf("Result: %+v\n", result)
		

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Scan completed successfully",
		})
	}
}