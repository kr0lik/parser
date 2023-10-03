package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
}

type ClientOptions struct {
	Host string
	User string
	Pass string
}

func NewClient(opts *ClientOptions, ctx context.Context) (*Client, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", opts.User, opts.Pass, opts.Host))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
