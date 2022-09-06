package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/riskiamad/go-pokemon/configs"
	pokemonController "github.com/riskiamad/go-pokemon/src/pokemon/controller"
	userController "github.com/riskiamad/go-pokemon/src/user/controller"
	userRepository "github.com/riskiamad/go-pokemon/src/user/repository"
	userService "github.com/riskiamad/go-pokemon/src/user/service"
)

var (
	env         = configs.Config
	mongoClient = configs.Client
	validate    = validator.New()
)

func main() {

	app := fiber.New()

	// set default cors config
	app.Use(cors.New())

	//routes group prefix
	v1 := app.Group("/api/v1")

	//dependencies
	userPokemonRepository := userRepository.NewUserPokemonRepository(mongoClient)
	userPokemonService := userService.NewUserPokemonService(userPokemonRepository, validate)
	userController.NewUserPokemonController(v1, userPokemonService)

	pokemonController.NewPokemonController(v1)

	log.Fatal(app.Listen(":" + env.Port))
}
