package handler

import (
	"github.com/ahmadirfaan/plantation-engine/usecase"
)

type Server struct {
	EstateUseCase      usecase.EstateUseCase
	TreeUseCase        usecase.TreeUseCase
	StatsEstateUseCase usecase.StatsEstateUseCase
}

type NewServerOptions struct {
	EstateUseCase      usecase.EstateUseCase
	TreeUseCase        usecase.TreeUseCase
	StatsEstateUseCase usecase.StatsEstateUseCase
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		EstateUseCase:      opts.EstateUseCase,
		TreeUseCase:        opts.TreeUseCase,
		StatsEstateUseCase: opts.StatsEstateUseCase,
	}
}
