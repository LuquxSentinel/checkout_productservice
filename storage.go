package main

import (
	"context"
	"log"
	"time"

	"github.com/luqu/productservice/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage interface {
	GetProduct(ctx context.Context, productID string) (*types.Product, error)
	GetAllProducts(ctx context.Context) ([]*types.Product, error)
}

type MongoStorage struct {
	productCollection *mongo.Collection
}

func InitDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoStorage(productCollection *mongo.Collection) *MongoStorage {
	return &MongoStorage{
		productCollection: productCollection,
	}
}

func (s *MongoStorage) GetProduct(ctx context.Context, productID string) (*types.Product, error) {
	filter := primitive.D{primitive.E{Key: "id", Value: productID}}

	product := new(types.Product)

	err := s.productCollection.FindOne(ctx, filter).Decode(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *MongoStorage) GetAllProducts(ctx context.Context) ([]*types.Product, error) {
	products := make([]*types.Product, 0)

	results, err := s.productCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = results.All(ctx, &products)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	results.Close(ctx)

	return products, nil
}
