package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"parser/internal/domain/wildberries/entity"
	"parser/internal/infrastructure/mongodb"
	"time"
)

type Category struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCategory(client *mongodb.Client, ctx context.Context) *Category {
	return &Category{
		collection: client.Database(databaseName).Collection("category"),
		ctx:        ctx,
	}
}

func (c *Category) IterateActive() <-chan *entity.Category {
	cur, err := c.collection.Find(c.ctx, bson.D{{"active", true}})
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *entity.Category)

	go func() {
		defer cur.Close(c.ctx)
		defer close(ch)

		for cur.Next(c.ctx) {
			category := new(entity.Category)

			if err := cur.Decode(category); err != nil {
				log.Fatal(err)
			}

			ch <- category
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return ch
}

func (c *Category) Upsert(category *entity.Category) error {
	filter := bson.M{"_id": category.ID}

	res := c.collection.FindOne(c.ctx, filter)
	if mongo.ErrNoDocuments == res.Err() {
		_, err := c.collection.InsertOne(c.ctx, category)
		return err
	}

	_, err := c.collection.ReplaceOne(c.ctx, filter, category)
	return err
}

func (c *Category) DisableOld(time time.Time) error {
	_, err := c.collection.UpdateMany(
		c.ctx,
		bson.M{"checked": bson.M{"$lt": primitive.NewDateTimeFromTime(time)}},
		bson.D{{"$set", bson.D{{"active", false}}}},
	)

	return err
}
