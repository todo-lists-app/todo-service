package service

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/keloran/go-healthcheck"
	"github.com/todo-lists-app/todo-service/internal/todo"
	"net"
	"net/http"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/todo/v1"
	"github.com/todo-lists-app/todo-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	*config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		Config: cfg,
	}
}

func (s *Service) Start() error {
	errChan := make(chan error)
	go startHTTP(s.Config, errChan)
	go startGRPC(s.Config, errChan)

	return nil
}

func startGRPC(cfg *config.Config, errChan chan error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Local.GRPCPort))
	if err != nil {
		errChan <- err
		return
	}

	gs := grpc.NewServer()
	reflection.Register(gs)
	pb.RegisterTodoServiceServer(gs, &todo.Server{
		Config: cfg,
	})
	logs.Local().Infof("starting grpc on port %d", cfg.Local.GRPCPort)
	if err := gs.Serve(lis); err != nil {
		errChan <- err
	}
}

func startHTTP(cfg *config.Config, errChan chan error) {
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://api.todo-list.app",
		"https://todo-list.app",
		"https://beta.todo-list.app",
	}
	if cfg.Local.Development {
		allowedOrigins = append(allowedOrigins, "http://*")
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET",
		},
	}))
	r.Get("/health", healthcheck.HTTP)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Local.HTTPPort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	logs.Local().Infof("starting http on port %d", cfg.Local.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		errChan <- err
	}
}
