package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"parser/internal/domain/ozon/entity"
	"parser/internal/infrastructure/mongodb"
)

type Product struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProduct(client *mongodb.Client, ctx context.Context) *Product {
	return &Product{
		collection: client.Database(databaseName).Collection("product"),
		ctx:        ctx,
	}
}

func (p *Product) Get(uuid string) (*entity.Product, error) {
	product := new(entity.Product)

	if err := p.collection.FindOne(p.ctx, bson.M{"_id": uuid}).Decode(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Upsert(product *entity.Product) error {
	filter := bson.M{"_id": product.ID}

	res := p.collection.FindOne(p.ctx, filter)
	if mongo.ErrNoDocuments == res.Err() {
		_, err := p.collection.InsertOne(p.ctx, product)
		return err
	}

	_, err := p.collection.ReplaceOne(p.ctx, filter, product)
	return err
}
