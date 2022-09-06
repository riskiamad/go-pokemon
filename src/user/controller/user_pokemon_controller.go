package controller

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/riskiamad/go-pokemon/dtos"
	"github.com/riskiamad/go-pokemon/models"
	"github.com/riskiamad/go-pokemon/utils"
)

var (
	responseFormat = utils.ResponseFormat{}
)

type userPokemonController struct {
	userPokemonService models.UserPokemonService
}

func NewUserPokemonController(route fiber.Router, service models.UserPokemonService) {
	controller := &userPokemonController{
		userPokemonService: service,
	}

	route.Post("/user/pokemon", controller.AddPokemon)
	route.Get("/user/pokemon", controller.GetPokemons)
	route.Get("/user/pokemon/:id", controller.GetPokemonById)
	route.Delete("/user/pokemon/:id/release", controller.ReleasePokemon)
	route.Patch("/user/pokemon/:id/rename", controller.RenamePokemon)
}

// AddPokemon: add pokemon to user pokemon list
func (c *userPokemonController) AddPokemon(ctx *fiber.Ctx) (err error) {
	var requestBody dtos.AddPokemon

	err = ctx.BodyParser(&requestBody)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(responseFormat.SetError(err))
	}

	result, err := c.userPokemonService.Store(ctx.Context(), &requestBody)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
	}

	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&result, nil))
}

func (c *userPokemonController) GetPokemons(ctx *fiber.Ctx) (err error) {
	var (
		queryParams = &dtos.QueryParams{
			Page:    1,
			PerPage: 10,
			OrderBy: "_id",
		}
		page    = ctx.Query("page")
		perpage = ctx.Query("perpage")
		orderby = ctx.Query("orderby")
	)

	if page != "" {
		convertedPage, err := strconv.Atoi(page)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
		}

		queryParams.Page = int64(convertedPage)
	}

	if perpage != "" {
		convertedPerPage, err := strconv.Atoi(perpage)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
		}

		queryParams.PerPage = int64(convertedPerPage)
	}

	if orderby != "" {
		queryParams.OrderBy = orderby
	}

	result, count, err := c.userPokemonService.GetPokemons(ctx.Context(), queryParams)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
	}

	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&result, count))
}

func (c *userPokemonController) GetPokemonById(ctx *fiber.Ctx) (err error) {
	var id = ctx.Params("id")

	result, err := c.userPokemonService.GetPokemonById(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
	}

	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&result, nil))
}

func (c *userPokemonController) ReleasePokemon(ctx *fiber.Ctx) (err error) {
	var id = ctx.Params("id")

	_, isReleased, err := c.userPokemonService.DeletePokemonById(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
	}

	data := map[string]bool{"is_released": isReleased}
	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&data, nil))
}

func (c *userPokemonController) RenamePokemon(ctx *fiber.Ctx) (err error) {
	var id = ctx.Params("id")

	result, err := c.userPokemonService.UpdatePokemonName(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(responseFormat.SetError(err))
	}

	return ctx.Status(http.StatusOK).JSON(responseFormat.SetData(&result, nil))
}
