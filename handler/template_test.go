package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Template_WhenBindRequestIsError(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`invalid-body`)

	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return nil },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{}, nil
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_HandlerTemplate_ValidationError(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`{"width":10,"length":100000}`)

	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return errors.New("validation failed") },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{}, nil
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_HandlerTemplate_HandleError(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`{"width":10,"length":100000}`)
	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return nil },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{}, errors.New("503|something went wrong")
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func Test_HandlerTemplate_HandleErrorMessage_SplitError(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`{"width":10,"length":100000}`)
	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return nil },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{}, errors.New("503|error format|multiple split")
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_HandlerTemplate_HandleErrorMessage_NoSplitMessage(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`{"width":10,"length":100000}`)
	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return nil },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{}, errors.New("timeout")
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_HandlerTemplate_SuccessProcess(t *testing.T) {
	rec, ctx := createRecorderAndContextEcho(`{"width":10,"length":100000}`)
	err := HandlerTemplate[generated.CreateEstateRequest, generated.CreatedResponse](
		ctx,
		func(req generated.CreateEstateRequest) error { return nil },
		func(req generated.CreateEstateRequest) (generated.CreatedResponse, error) {
			return generated.CreatedResponse{Id: uuid.New().String()}, nil
		},
		func(ctx echo.Context, res generated.CreatedResponse) error {
			return ctx.JSON(http.StatusOK, res)
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func createRecorderAndContextEcho(requestBody string) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return rec, ctx
}
