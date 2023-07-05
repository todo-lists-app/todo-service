package todo

import (
	"context"
	"errors"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-service/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	config.Config
	context.Context
	UserID string
}

type List struct {
	UserID string `bson:"userid" json:"userid"`
	Data   string `bson:"data" json:"data"`
	IV     string `bson:"iv" json:"iv"`
}

func NewTodoService(ctx context.Context, cfg config.Config, userId string) *Service {
	return &Service{
		Config:  cfg,
		Context: ctx,
		UserID:  userId,
	}
}

func (t *Service) GetTodoList() (*List, error) {
	client, err := config.GetMongoClient(t.Context, t.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	tl := List{}
	if err := client.Database(t.Config.Mongo.Database).Collection(t.Config.Mongo.Collections.List).FindOne(t.Context, &bson.M{
		"userid": t.UserID,
	}).Decode(&tl); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return &tl, logs.Errorf("error finding list: %v", err)
		}
		return &tl, nil
	}

	return &tl, nil
}

func (t *Service) UpdateTodoList(data, iv string) (*List, error) {
	client, err := config.GetMongoClient(t.Context, t.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	if _, err := client.Database(t.Config.Mongo.Database).Collection(t.Config.Mongo.Collections.List).UpdateOne(t.Context,
		bson.D{{"userid", t.UserID}},
		bson.D{{"$set", bson.M{
			"data": data,
			"iv":   iv,
		}}}); err != nil {
		return nil, logs.Errorf("error updating list: %v", err)
	}

	return &List{
		UserID: t.UserID,
		Data:   data,
		IV:     iv,
	}, nil
}

func (t *Service) DeleteTodoList() error {
	return nil
}

func (t *Service) CreateTodoList(data, iv string) (*List, error) {
	client, err := config.GetMongoClient(t.Context, t.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()
	if _, err := client.Database(t.Config.Mongo.Database).Collection(t.Config.Mongo.Collections.List).InsertOne(t.Context, &List{
		UserID: t.UserID,
		Data:   data,
		IV:     iv,
	}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &List{
				UserID: t.UserID,
				Data:   data,
				IV:     iv,
			}, nil
		}
		return nil, logs.Errorf("error creating list: %v", err)
	}

	return &List{
		UserID: t.UserID,
		Data:   data,
		IV:     iv,
	}, nil
}
