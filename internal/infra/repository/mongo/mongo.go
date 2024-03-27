package mongo

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	return &MongoRepository{
		client: client,
	}, nil
}

func (r *MongoRepository) Ping(c context.Context) error {
	return r.client.Ping(c, readpref.Primary())
}

func (r *MongoRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {

}

func (r *MongoRepository) Create(c context.Context, edge domain.Edge) error {
	col := r.client.Database(viper.GetString("db")).Collection("edges")
	_, err := col.InsertOne(
		c,
		edge,
	)
	return err
}

func (r *MongoRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {

}

func (r *MongoRepository) ClearAll(c context.Context) error {

}
