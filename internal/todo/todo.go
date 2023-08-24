package todo

import (
	"context"
	"errors"
	"github.com/bugfixes/go-bugfixes/logs"
	mungo "github.com/keloran/go-config/mongo"
	"github.com/todo-lists-app/todo-service/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOperations interface {
	GetMongoClient(ctx context.Context, config mungo.Mongo) error
	Disconnect(ctx context.Context) error
	InsertOne(ctx context.Context, document interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error)
	FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult
}

type RealMongoOperations struct {
	Client     *mongo.Client
	Collection string
	Database   string
}

func (r *RealMongoOperations) GetMongoClient(ctx context.Context, config mungo.Mongo) error {
	client, err := mungo.GetMongoClient(ctx, config)
	if err != nil {
		return logs.Errorf("error getting mongo client: %v", err)
	}
	r.Client = client
	return nil
}
func (r *RealMongoOperations) Disconnect(ctx context.Context) error {
	return r.Client.Disconnect(ctx)
}
func (r *RealMongoOperations) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	return r.Client.Database(r.Database).Collection(r.Collection).InsertOne(ctx, document)
}
func (r *RealMongoOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error) {
	return r.Client.Database(r.Database).Collection(r.Collection).UpdateOne(ctx, filter, update)
}
func (r *RealMongoOperations) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	return r.Client.Database(r.Database).Collection(r.Collection).FindOne(ctx, filter)
}

type Service struct {
	config.Config
	context.Context
	UserID string

	MongoOps MongoOperations
}

type List struct {
	UserID string `bson:"userid" json:"userid"`
	Data   string `bson:"data" json:"data"`
	IV     string `bson:"iv" json:"iv"`
}

func NewTodoService(ctx context.Context, cfg config.Config, userId string, ops MongoOperations) *Service {
	return &Service{
		Config:   cfg,
		Context:  ctx,
		UserID:   userId,
		MongoOps: ops,
	}
}

func (t *Service) GetTodoList() (*List, error) {
	if err := t.MongoOps.GetMongoClient(t.Context, t.Config.Mongo); err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := t.MongoOps.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	tl := List{}
	if err := t.MongoOps.FindOne(t.Context, &bson.M{
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
	if err := t.MongoOps.GetMongoClient(t.Context, t.Config.Mongo); err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := t.MongoOps.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	if _, err := t.MongoOps.UpdateOne(t.Context,
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
	if err := t.MongoOps.GetMongoClient(t.Context, t.Config.Mongo); err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := t.MongoOps.Disconnect(t.Context); err != nil {
			_ = logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()
	if _, err := t.MongoOps.InsertOne(t.Context, &List{
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
