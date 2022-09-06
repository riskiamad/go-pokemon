package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
	"github.com/mtslzr/pokeapi-go"
	"github.com/riskiamad/go-pokemon/dtos"
	"github.com/riskiamad/go-pokemon/models"
	"github.com/riskiamad/go-pokemon/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userPokemonService struct {
	userPokemonRepository models.UserPokemonRepository
	validate              *validator.Validate
}

func NewUserPokemonService(repository models.UserPokemonRepository, validate *validator.Validate) models.UserPokemonService {
	return &userPokemonService{
		userPokemonRepository: repository,
		validate:              validate,
	}
}

func (s *userPokemonService) Store(ctx context.Context, req *dtos.AddPokemon) (*models.UserPokemon, error) {
	err := s.validate.Struct(req)
	if err != nil {
		return nil, err
	}

	// convert pokemonID to string
	strPokemonID := strconv.FormatInt(req.PokemonId, 10)

	// call pokemon struct with pokeapi wrapper
	pokemon, err := pokeapi.Pokemon(strPokemonID)
	if err != nil {
		return nil, err
	}

	// convert pokemon struct to map
	pokemonMap := structs.Map(pokemon)

	userPokemon := &models.UserPokemon{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Pokemon:     pokemonMap,
		RenameCount: 0,
		CreatedAt:   time.Now(),
	}

	res, err := s.userPokemonRepository.Insert(ctx, userPokemon)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userPokemonService) GetPokemons(ctx context.Context, queryParams *dtos.QueryParams) ([]*models.UserPokemon, int64, error) {
	res, count, err := s.userPokemonRepository.Find(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	return res, count, err
}

func (s *userPokemonService) GetPokemonById(ctx context.Context, id string) (*models.UserPokemon, error) {
	res, err := s.userPokemonRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userPokemonService) DeletePokemonById(ctx context.Context, id string) (*models.UserPokemon, bool, error) {
	userPokemon, err := s.userPokemonRepository.FindById(ctx, id)
	if err != nil {
		return nil, false, err
	}

	// get random number between 1-20
	randNum := utils.RandomNum(1, 20)

	// check if the random number is prime
	isPrime := utils.CheckIsPrime(randNum)

	if !isPrime {
		return userPokemon, false, nil
	}

	err = s.userPokemonRepository.DeleteById(ctx, id)
	if err != nil {
		return nil, false, nil
	}

	return userPokemon, true, err
}

func (s *userPokemonService) UpdatePokemonName(ctx context.Context, id string) (m *models.UserPokemon, err error) {
	userPokemon, err := s.userPokemonRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// get fibonacci sequence number
	sequence := utils.FibonacciSequence(userPokemon.RenameCount)
	suffix := strconv.Itoa(sequence)

	name := userPokemon.Name

	if userPokemon.RenameCount > 0 {
		// get index of suffix name, is suffix not found return -1
		suffixIndex := strings.LastIndexAny(userPokemon.Name, "-")
		name = userPokemon.Name[:suffixIndex]
	}

	userPokemon.Name = name + "-" + suffix
	userPokemon.RenameCount++
	userPokemon.UpdatedAt = time.Now()

	m, err = s.userPokemonRepository.UpdateNameByID(ctx, userPokemon, id)
	if err != nil {
		return nil, err
	}

	return m, err
}
