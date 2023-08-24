package todo

import (
	"context"
	"errors"
	ConfigBuilder "github.com/keloran/go-config"
	mungo "github.com/keloran/go-config/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/todo-lists-app/todo-service/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type MockMongoOperations struct {
	shouldError bool
	exists      bool
}

func (m *MockMongoOperations) GetMongoClient(ctx context.Context, config mungo.Mongo) error {
	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}

func (m *MockMongoOperations) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockMongoOperations) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	return nil, nil
}

func (m *MockMongoOperations) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	return nil, nil
}

func (m *MockMongoOperations) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	doc := bson.D{{"userid", "test"}, {"time", time.Now()}}

	if m.shouldError {
		// This should be adjusted based on how you handle errors in the FindOne method.
		return mongo.NewSingleResultFromDocument(nil, errors.New("mock error"), bson.DefaultRegistry)
	}
	if m.exists {
		// Mocked result
		return mongo.NewSingleResultFromDocument(doc, nil, bson.DefaultRegistry)
	}
	// Simulate a "not found" scenario

	return mongo.NewSingleResultFromDocument(doc, mongo.ErrNoDocuments, bson.DefaultRegistry)
}

func TestService_GetTodoList(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Config: ConfigBuilder.Config{
			Mongo: mungo.Mongo{
				Database:    "testDB",
				Collections: map[string]string{"todo": "testCollection"},
			},
		},
	}

	t.Run("GetTodoList", func(t *testing.T) {
		todo := NewTodoService(ctx, cfg, "test", &MockMongoOperations{})
		_, err := todo.GetTodoList()
		assert.Nil(t, err)
	})
}
