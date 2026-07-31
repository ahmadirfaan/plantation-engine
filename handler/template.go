package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/labstack/echo/v4"
)

func HandlerTemplate[Req any, Res any](
	ctx echo.Context,
	doValidate func(req Req) error,
	doHandle func(req Req) (Res, error),
	respond func(ctx echo.Context, res Res) error,
) error {
	var req Req

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	err := doValidate(req)
	if err != nil {
		log.Println("validation error:", err)
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	res, err := doHandle(req)
	if err != nil {
		log.Println("doHandle error:", err)
		return handleErrorFromDoHandle(ctx, err)
	}

	return respond(ctx, res)
}

func handleErrorFromDoHandle(ctx echo.Context, err error) error {
	errorMessage := err.Error()
	if strings.Contains(errorMessage, "|") {
		errorPart := strings.Split(errorMessage, "|")
		if len(errorPart) == 2 {
			statusCode := errorPart[0]
			message := errorPart[1]
			code, _ := strconv.Atoi(statusCode)
			return ctx.JSON(code, generated.ErrorResponse{Message: message})
		} else {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	} else {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}
}
