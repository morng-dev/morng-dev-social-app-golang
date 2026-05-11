package validation

import (
	"Server/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ValidatorUser = validator.New()

func ValidateUser(c *fiber.Ctx) error {
	var errors []*model.IError
	var body model.UserModel

	if err := c.BodyParser(&body); err != nil {
		return err
	}
	err := ValidatorUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el model.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}
