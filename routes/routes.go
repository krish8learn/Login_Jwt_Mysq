package routes

import(
	"../controllers"
	"github.com/gofiber/fiber"
)

func Setup(app *fiber.App){
	app.Post("/api/register",controllers.Regisiter)
	app.Post("/api/login",controllers.Login)
	app.Get("/api/user",controllers.User)
	app.Post("/api/logout",controllers.LogOut)
}
