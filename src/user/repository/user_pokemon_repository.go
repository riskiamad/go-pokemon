package repository

import (
	"context"

	"github.com/riskiamad/go-pokemon/configs"
	"github.com/riskiamad/go-pokemon/dtos"
	"github.com/riskiamad/go-pokemon/models"
	"github.com/riskiamad/go-pokemon/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type userPokemonRepository struct {
	client *mongo.Client
}

func NewUserPokemonRepository(client *mongo.Client) models.UserPokemonRepository {
	return &userPokemonRepository{
		client: client,
	}
}

func (r *userPokemonRepository) Insert(ctx context.Context, userPokemon *models.UserPokemon) (m *models.UserPokemon, err error) {
	var userPokemonCollection = configs.GetCollection(r.client, "user_pokemon")
	wc := writeconcern.New(writeconcern.WMajority())
	txOpt := options.Transaction().SetWriteConcern(wc)

	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		res, err := userPokemonCollection.InsertOne(
			sessionContext,
			userPokemon,
		)
		if err != nil {
			return nil, err
		}

		return res, err
	}

	_, err = session.WithTransaction(ctx, callback, txOpt)
	if err != nil {
		return nil, err
	}

	return userPokemon, nil
}

func (r *userPokemonRepository) Find(ctx context.Context, params *dtos.QueryParams) (mx []*models.UserPokemon, count int64, err error) {
	var (
		userPokemonCollection = configs.GetCollection(r.client, "user_pokemon")
		filter                map[string]interface{}
		opts                  *options.FindOptions
	)

	orderby := utils.SetSort(params.OrderBy)

	skip := (params.Page * params.PerPage) - params.PerPage
	opts = options.MergeFindOptions(
		options.Find().SetLimit(params.PerPage),
		options.Find().SetSkip(skip),
		options.Find().SetSort(orderby),
	)
	cursor, err := userPokemonCollection.Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, err
	}

	err = cursor.All(ctx, &mx)
	if err != nil {
		return nil, 0, err
	}

	count, err = userPokemonCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return mx, count, nil
}

func (r *userPokemonRepository) FindById(ctx context.Context, id string) (m *models.UserPokemon, err error) {
	var userPokemonCollection = configs.GetCollection(r.client, "user_pokemon")

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = userPokemonCollection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *userPokemonRepository) DeleteById(ctx context.Context, id string) (err error) {
	var userPokemonCollection = configs.GetCollection(r.client, "user_pokemon")
	wc := writeconcern.New(writeconcern.WMajority())
	txOpt := options.Transaction().SetWriteConcern(wc)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		res, err := userPokemonCollection.DeleteOne(
			sessionContext,
			bson.M{"_id": idHex},
		)
		if err != nil {
			return nil, err
		}

		return res, err
	}

	_, err = session.WithTransaction(ctx, callback, txOpt)
	if err != nil {
		return err
	}

	return nil
}

func (r *userPokemonRepository) UpdateNameByID(ctx context.Context, userPokemon *models.UserPokemon, id string) (m *models.UserPokemon, err error) {
	var userPokemonCollection = configs.GetCollection(r.client, "user_pokemon")
	wc := writeconcern.New(writeconcern.WMajority())
	txOpt := options.Transaction().SetWriteConcern(wc)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return userPokemon, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": userPokemon}

	session, err := r.client.StartSession()
	if err != nil {
		return userPokemon, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		res, err := userPokemonCollection.UpdateOne(
			sessionContext,
			filter,
			update,
		)
		if err != nil {
			return nil, err
		}

		return res, err
	}

	_, err = session.WithTransaction(ctx, callback, txOpt)
	if err != nil {
		return userPokemon, err
	}

	err = userPokemonCollection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
