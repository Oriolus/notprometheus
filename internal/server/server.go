package server

import (
	"github.com/oriolus/notprometheus/internal/server/storage"
)

type Server struct {
	storage storage.Storage
}

func NewServer(storage storage.Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) Storage() storage.Storage {
	return s.storage
}
