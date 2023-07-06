package todo

import (
	"context"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	pb "github.com/todo-lists-app/protobufs/generated/todo/v1"
	"github.com/todo-lists-app/todo-service/internal/config"
)

type Server struct {
	pb.UnimplementedTodoServiceServer
	*config.Config
}

func (s *Server) Get(ctx context.Context, r *pb.TodoGetRequest) (*pb.TodoRetrieveResponse, error) {
	t := NewTodoService(ctx, *s.Config, r.UserId)
	list, err := t.GetTodoList()
	if err != nil {
		return &pb.TodoRetrieveResponse{
			UserId: r.UserId,
			Data:   "",
			Iv:     "",
			Status: pointerutil.StringPtr(err.Error()),
		}, err
	}

	return &pb.TodoRetrieveResponse{
		UserId: r.UserId,
		Data:   list.Data,
		Iv:     list.IV,
	}, nil
}

func (s *Server) Insert(ctx context.Context, r *pb.TodoInjectRequest) (*pb.TodoRetrieveResponse, error) {
	t := NewTodoService(ctx, *s.Config, r.UserId)
	list, err := t.GetTodoList()
	if err != nil {
		return &pb.TodoRetrieveResponse{
			UserId: r.UserId,
			Data:   "",
			Iv:     "",
			Status: pointerutil.StringPtr(err.Error()),
		}, nil
	}

	// data is empty, create a new list
	list, err = t.CreateTodoList(r.Data, r.Iv)
	if err != nil {
		return &pb.TodoRetrieveResponse{
			UserId: r.UserId,
			Data:   "",
			Iv:     "",
			Status: pointerutil.StringPtr(err.Error()),
		}, nil
	}
	return &pb.TodoRetrieveResponse{
		UserId: r.UserId,
		Data:   list.Data,
		Iv:     list.IV,
	}, nil
}

func (s *Server) Update(ctx context.Context, r *pb.TodoInjectRequest) (*pb.TodoRetrieveResponse, error) {
	t := NewTodoService(ctx, *s.Config, r.UserId)
	list, err := t.UpdateTodoList(r.Data, r.Iv)
	if err != nil {
		return &pb.TodoRetrieveResponse{
			UserId: r.UserId,
			Data:   "",
			Iv:     "",
			Status: pointerutil.StringPtr(err.Error()),
		}, err
	}

	return &pb.TodoRetrieveResponse{
		UserId: r.UserId,
		Data:   list.Data,
		Iv:     list.IV,
	}, nil
}
