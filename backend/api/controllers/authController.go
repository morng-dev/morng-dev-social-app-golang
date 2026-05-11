package controllers

import (
	"Server/database"
	_ "Server/docs"
	"Server/model"
	"Server/util"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Register
// @Summary Register a new user
// @Description Register a new user by providing email, password, first name, and last name
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body model.CreateUser true "User registration details"
// @Success 201 {object} model.UserModel "User registered successfully"
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	var body model.CreateUser
	var userSchema = database.DB.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	CheckUser := userSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&body)

	if CheckUser == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "User all ready exist,!!!",
		})
	}
	hashPassword, err := util.HashPassword(body.Password)
	if err != nil {
		return err
	}
	newUser := model.UserModel{
		Name:       body.FirstName + " " + body.LastName,
		Email:      body.Email,
		Password:   hashPassword,
		Followers:  make([]string, 0),
		Followeing: make([]string, 0),
	}

	result, err := userSchema.InsertOne(ctx, newUser)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	//get new user
	var createdUser *model.UserModel
	query := bson.M{"_id": result.InsertedID}
	userSchema.FindOne(ctx, query).Decode(&createdUser)
	// createToken
	token, _ := util.GenerateJWT(createdUser.ID.Hex())

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": createdUser,
		"token":  token,
	})
}

func Login(c *fiber.Ctx) error {
	var userSchema = database.DB.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var body model.LoginUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "invalid requrest body",
			"detail": err.Error(),
		})
	}
	var user model.UserModel
	CheckEmail := userSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if CheckEmail != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "email or password invalid",
			"success": false,
		})
	}
	if !util.CheckPassword(user.Password, body.Password) {
		return errors.New("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}
	token, _ := util.GenerateJWT(user.ID.Hex())
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": user,
		"token":  token,
	})
}
