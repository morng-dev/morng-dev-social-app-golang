package controllers

import (
	"Server/database"
	_ "Server/docs"
	"Server/model"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserByID
// @Summary GetUserByID
// @Description GetUserByID detail By ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Failure 400 {object} map[string]interface{}
// @Router /user/:id [Get]
func GetUserByID(c *fiber.Ctx) error {
	var userSchema = database.DB.Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	user := model.UserModel{}

	objID, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}
	userResult := userSchema.FindOne(ctx, bson.M{"_id": objID})
	if userResult.Err() != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
		})
	}

	userResult.Decode(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user,
		"posts": []string{},
	})
}

func Update(c *fiber.Ctx) error {
	var userSchema = database.DB.Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	extUid, ok := c.Locals("userId").(string)
	fmt.Println(extUid)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid userID type",
		})
	}
	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "you are not authorization update this profile",
		})
	}
	userid, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}
	user := model.UpdateUser{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	update := bson.M{"name": user.Name, "imageUrl": user.ImageUrl, "bio": user.Bio}

	result, err := userSchema.UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "cannot update user",
		})
	}
	updateUser := model.UserModel{}

	if result.MatchedCount == 1 {
		err := userSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&updateUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "cannot update user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})

}
