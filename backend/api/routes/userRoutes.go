package routes

import (
	"Server/controllers"
	"Server/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	//auth
	app.Get("/user/getuser/:id", controllers.GetUserByID)
	//getslug
	//update
	app.Patch("/user/update/:id", middleware.AuthMiddleware, controllers.Update)

}
