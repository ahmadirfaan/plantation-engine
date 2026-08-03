package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ahmadirfaan/plantation-engine/validator"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) PostEstate(ctx echo.Context) error {
	return HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,

		func(req generated.CreateEstateRequest) error {
			return validator.ValidateRequestCreateEstate(req)
		},

		// Business logic
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			id, err := s.EstateUseCase.CreateEstate(ctx.Request().Context(), req)
			if err != nil {
				return generated.CreatedResponse{}, errors.New("500|" + err.Error())
			}
			return generated.CreatedResponse{Id: *id}, nil
		},

		// Response encoder
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusCreated, res)
		},
	)
}

func (s *Server) AddTreeToEstate(ctx echo.Context, estateId string) error {
	return HandlerTemplate[generated.CreateTreeRequest, generated.CreatedResponse](ctx,
		func(req generated.CreateTreeRequest) error {

			return validator.ValidateRequestCreateTree(req, estateId)
		},
		func(req generated.CreateTreeRequest) (generated.CreatedResponse, error) {
			id, err := s.TreeUseCase.AddTreeToEstate(ctx.Request().Context(), estateId, req)
			if err != nil {
				return generated.CreatedResponse{}, err
			}
			return generated.CreatedResponse{Id: *id}, nil
		}, func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusCreated, res)
		})
}

func (s *Server) GetEstateSummary(ctx echo.Context, estateId string) error {
	return HandlerTemplate[any, generated.GetStatsEstateResponse](ctx,
		func(req any) error {

			return validator.ValidateEstateId(estateId)
		},
		func(req any) (generated.GetStatsEstateResponse, error) {
			summaryEstate, err := s.StatsEstateUseCase.GetStatsEstate(ctx.Request().Context(), estateId)
			if err != nil {
				return generated.GetStatsEstateResponse{}, err
			}
			return summaryEstate, nil
		}, func(ctx echo.Context, res generated.GetStatsEstateResponse) error {
			return ctx.JSON(http.StatusOK, res)
		})
}

func (s *Server) GetDistanceForDronePlan(ctx echo.Context, estateId string, params generated.GetDistanceForDronePlanParams) error {
	return HandlerTemplate[any, generated.GetDronePlanDistance](ctx,
		func(req any) error {

			return validator.ValidateEstateId(estateId)
		},
		func(req any) (generated.GetDronePlanDistance, error) {
			dronePlanDistance, err := s.StatsEstateUseCase.GetDronePlanDistance(ctx.Request().Context(), estateId, params)
			if err != nil {
				return generated.GetDronePlanDistance{}, err
			}
			return dronePlanDistance, nil
		}, func(ctx echo.Context, res generated.GetDronePlanDistance) error {
			return ctx.JSON(http.StatusOK, res)
		})
}
