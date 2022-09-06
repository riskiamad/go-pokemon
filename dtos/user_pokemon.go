package dtos

// AddPokemon: request body for add pokemon to user pokemon list
type AddPokemon struct {
	Name      string `json:"name" validate:"required"`
	PokemonId int64  `json:"pokemon_id" validate:"required"`
}
