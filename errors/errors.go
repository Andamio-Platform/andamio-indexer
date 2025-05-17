package errors

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// Return JSON error model
type Message struct {
	ID        interface{} `json:"ID,omitempty"`
	IsError   bool        `json:"IsError,omitempty"`
	Message   string
	Necessary []string `json:"Necessary,omitempty"`
}

type ServerError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

func FieldError(str ...string) *Message {
	return &Message{
		IsError:   true,
		Message:   "Necessary fields can't be empty",
		Necessary: str,
	}
}

func GeneralError(str string) *Message {
	return &Message{
		IsError: true,
		Message: str,
	}
}

func NoError(str string, id interface{}) *Message {
	return &Message{
		ID:      id,
		IsError: false,
		Message: str,
	}
}

func RouteError(c *fiber.Ctx) error {

	c.Status(fiber.StatusNotFound)
	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")

	return enc.Encode(&Message{
		IsError: true,
		Message: "THERE IS NO SUCH PAGE",
	})

}


func ServerErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Check if it's an fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Status(code)
	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")
	return enc.Encode(&ServerError{
		StatusCode: code,
		Error:      err.Error(),
	})
}

func BadRequestErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusBadRequest

	// Check if it's an fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Status(code)
	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")
	return enc.Encode(&ServerError{
		StatusCode: code,
		Error:      err.Error(),
	})
}
