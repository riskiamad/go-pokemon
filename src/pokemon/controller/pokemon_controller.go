package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/riskiamad/go-pokemon/utils"
)

type PokemonController struct{}

func NewPokemonController(route fiber.Router) {
	controller := &PokemonController{}

	route.Get("/pokemon/catch", controller.Catch)
}

func (c *PokemonController) Catch(ctx *fiber.Ctx) (err error) {
	var responseFormat utils.ResponseFormat
	isCaught := utils.RandomNum(0, 1) == 1
	result := map[string]bool{"is_caught": isCaught}
	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&result, nil))
}
