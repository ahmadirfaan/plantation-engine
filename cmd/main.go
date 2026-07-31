package main

import (
	"database/sql"
	"os"

	"github.com/ahmadirfaan/plantation-engine/usecase"
	_ "github.com/lib/pq"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/handler"
	"github.com/ahmadirfaan/plantation-engine/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := NewServerApp()
	e.Logger.Fatal(e.Start(":1323"))
}

func NewServerApp() *echo.Echo {
	dbDsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		panic(err)
	}

	estateRepository := repository.NewEstateRepository(db)
	blockRepository := repository.NewBlockRepository(db)
	treeRepository := repository.NewTreeRepository(db)
	statsEstateRepository := repository.NewStatsEstateRepository(db)
	statsEstateUseCase := usecase.NewStatsEstateUseCase(statsEstateRepository, estateRepository)
	estateUseCase := usecase.NewEstateUseCase(estateRepository)
	treeUseCase := usecase.NewTreeUseCase(estateRepository, treeRepository, blockRepository, statsEstateUseCase)

	opts := handler.NewServerOptions{
		EstateUseCase:      estateUseCase,
		TreeUseCase:        treeUseCase,
		StatsEstateUseCase: statsEstateUseCase,
	}

	handlerServer := handler.NewServer(opts)
	e := echo.New()
	generated.RegisterHandlers(e, handlerServer)
	e.Use(middleware.Logger())
	return e
}
