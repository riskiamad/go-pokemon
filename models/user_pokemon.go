package models

import (
	"context"
	"time"

	"github.com/riskiamad/go-pokemon/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserPokemon: struct hold data of pokemon which owned by user
type UserPokemon struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	// User *User `bson:"user"`
	// PokemonID   int64                  `bson:"pokemon_id"`
	Pokemon     map[string]interface{} `bson:"pokemon" json:"pokemon"`
	RenameCount int                    `bson:"rename_count" json:"rename_count"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

// UserPokemonService: represent the user pokemon service
type UserPokemonService interface {
	Store(ctx context.Context, request *dtos.AddPokemon) (*UserPokemon, error)
	GetPokemons(ctx context.Context, queryParams *dtos.QueryParams) ([]*UserPokemon, int64, error)
	GetPokemonById(ctx context.Context, id string) (*UserPokemon, error)
	DeletePokemonById(ctx context.Context, id string) (*UserPokemon, bool, error)
	UpdatePokemonName(ctx context.Context, id string) (*UserPokemon, error)
}

// UserPokemonRepository: represent the user pokemon repository
type UserPokemonRepository interface {
	Insert(ctx context.Context, userPokemon *UserPokemon) (*UserPokemon, error)
	Find(ctx context.Context, queryParams *dtos.QueryParams) ([]*UserPokemon, int64, error)
	FindById(ctx context.Context, id string) (*UserPokemon, error)
	DeleteById(ctx context.Context, id string) error
	UpdateNameByID(ctx context.Context, userPokemon *UserPokemon, id string) (*UserPokemon, error)
}
