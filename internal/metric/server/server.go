package server

import "github.com/oriolus/notprometheus/internal/metric/storage"

type Server struct {
	storage storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	return &Server{storage: storage}
}
