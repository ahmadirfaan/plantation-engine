package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/repository"
	"github.com/ahmadirfaan/plantation-engine/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db         *sql.DB
	estateRepo repository.EstateRepository
	server     *Server
)

type mockEstateUseCaseError struct{}

func (m *mockEstateUseCaseError) CreateEstate(ctx context.Context, req generated.CreateEstateRequest) (*string, error) {
	return nil, errors.New("mock create estate error")
}

type mockTreeUseCaseError struct{}

func (m mockTreeUseCaseError) AddTreeToEstate(ctx context.Context, estateId string, request generated.CreateTreeRequest) (*string, error) {
	return nil, errors.New("mock add tree error")
}

func assertBadRequest(t *testing.T, rec *httptest.ResponseRecorder, expectedMessage string) {
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp generated.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, resp.Message)
}

func createContextAndRecordHttp(body string, url string) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)
	return rec, ctx
}

func TestMain(m *testing.M) {
	var err error
	// Initialize Echo instance
	_ = echo.New()

	// Initialize database connection
	dsn := "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"
	log.Println("dsn:", dsn)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Create the repository and use case
	estateRepo = repository.NewEstateRepository(db)
	blockRepository := repository.NewBlockRepository(db)
	treeRepository := repository.NewTreeRepository(db)
	statsEstateRepository := repository.NewStatsEstateRepository(db)
	statsEstateUseCase := usecase.NewStatsEstateUseCase(statsEstateRepository, estateRepo)
	estateUseCase := usecase.NewEstateUseCase(estateRepo)
	treeUseCase := usecase.NewTreeUseCase(estateRepo, treeRepository, blockRepository, statsEstateUseCase)

	// Initialize the server with the use case
	server = NewServer(NewServerOptions{
		EstateUseCase:      estateUseCase,
		TreeUseCase:        treeUseCase,
		StatsEstateUseCase: statsEstateUseCase,
	})

	code := m.Run()
	_ = db.Close()
	os.Exit(code)

}

func TestPostEstateWidthExceedLimit(t *testing.T) {
	body := `{"width":100000,"length":100}`
	rec, ctx := createContextAndRecordHttp(body, "/estate")

	_ = server.PostEstate(ctx)

	assertBadRequest(t, rec, "width must be between 1 and 5000")

}

func TestPostEstateLengthExceedLimit(t *testing.T) {
	body := `{"width":10,"length":100000}`
	rec, ctx := createContextAndRecordHttp(body, "/estate")

	_ = server.PostEstate(ctx)

	assertBadRequest(t, rec, "length must be between 1 and 5000")
}

func TestGetHello(t *testing.T) {
	rec, ctx := createContextAndRecordHttp("", "/hello?123")
	_ = server.GetHello(ctx, generated.GetHelloParams{Id: 123})

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp generated.HelloResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello User 123", resp.Message)

}

func Test_PostEstate_CreateEstateError(t *testing.T) {
	body := `{"width":10,"length":1000}`
	rec, ctx := createContextAndRecordHttp(body, "/estate")

	server := &Server{
		EstateUseCase: &mockEstateUseCaseError{},
	}

	err := server.PostEstate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "mock create estate error")
}

func Test_PostEstate_CreateEstateSuccess(t *testing.T) {
	requireTestDB(t)
	body := `{"width":10,"length":1000}`
	resp := createEstate(t, body)
	assert.NotNil(t, resp.Id)
}

func requireTestDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}
	if err := db.Ping(); err != nil {
		t.Skipf("skipping E2E test, database not reachable: %v", err)
	}
}

func createEstate(t *testing.T, body string) generated.CreatedResponse {
	rec, ctx := createContextAndRecordHttp(body, "/estate")

	err := server.PostEstate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp generated.CreatedResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	return resp
}

func TestAddTreeToEstate_LongitudeExceedLimit(t *testing.T) {

	estateId := uuid.New().String()
	body := `{"x":100000,"y":6000,"height":20}`
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")
	_ = server.AddTreeToEstate(ctx, estateId)
	assertBadRequest(t, rec, "x must be between 1 and 5000")

}

func TestAddTreeToEstate_LatitudeExceedLimit(t *testing.T) {
	estateId := uuid.New().String()
	body := `{"x":100,"y":100000,"height":20}`
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")
	_ = server.AddTreeToEstate(ctx, estateId)
	assertBadRequest(t, rec, "y must be between 1 and 5000")
}

func TestAddTreeToEstate_HeightExceedLimit(t *testing.T) {
	estateId := uuid.New().String()
	body := `{"x":100,"y":100,"height":31}`
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")
	_ = server.AddTreeToEstate(ctx, estateId)
	assertBadRequest(t, rec, "height must be between 1 and 30")
}

func TestAddTreeToEstate_ErrorEstateId(t *testing.T) {
	newUUID, _ := uuid.NewUUID()
	estateId := newUUID.String()
	body := `{"x":100,"y":100,"height":31}`
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")
	_ = server.AddTreeToEstate(ctx, estateId)
	assertBadRequest(t, rec, "400|must valid estate id")
}

func TestAddTreeToEstate_InternalServerError(t *testing.T) {
	estateId := uuid.New().String()
	body := `{"x":10,"y":15,"height":15}`
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")

	server := &Server{
		TreeUseCase: &mockTreeUseCaseError{},
	}

	err := server.AddTreeToEstate(ctx, estateId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "mock add tree error")
}

func TestAddTreeToEstate_Success(t *testing.T) {
	requireTestDB(t)

	bodyPostEstate := `{"width":1000,"length":10}`
	createResponseSuccess := createEstate(t, bodyPostEstate)
	estateId := createResponseSuccess.Id
	body := `{"x":10,"y":15,"height":15}`
	resp := addTreeToEstate(t, body, estateId)
	assert.True(t, strings.TrimSpace(resp.Id) != "")

}

func addTreeToEstate(t *testing.T, body string, estateId string) generated.CreatedResponse {
	rec, ctx := createContextAndRecordHttp(body, "/estate/"+estateId+"/tree")

	err := server.AddTreeToEstate(ctx, estateId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp generated.CreatedResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotNil(t, resp)
	return resp
}

func TestGetSummaryEstate_ErrorEstateId(t *testing.T) {
	newUUID, _ := uuid.NewUUID()
	estateId := newUUID.String()
	rec, ctx := createContextAndRecordHttp("", "/estate/"+estateId+"/stats")
	_ = server.GetEstateSummary(ctx, estateId)
	assertBadRequest(t, rec, "400|must valid estate id")
}

func TestGetSummaryEstate_And_DistanceDronePlan_Success(t *testing.T) {
	requireTestDB(t)

	bodyCreateEstate := `{"width":1,"length":5}`
	respCreatedEstate := createEstate(t, bodyCreateEstate)
	estateId := respCreatedEstate.Id

	allRequestCreateTree := make([]string, 0)
	allRequestCreateTree = append(allRequestCreateTree, `{"x":2,"y":1,"height":5}`)
	allRequestCreateTree = append(allRequestCreateTree, `{"x":3,"y":1,"height":3}`)
	allRequestCreateTree = append(allRequestCreateTree, `{"x":4,"y":1,"height":4}`)

	for _, request := range allRequestCreateTree {
		addTreeToEstate(t, request, estateId)
	}

	waitForStatsToBeCalculated(t, estateId, 1*time.Second)
}

func TestGetSummaryEstate_InternalServerError(t *testing.T) {
	estateId := uuid.New().String()
	rec, ctx := createContextAndRecordHttp("", "/estate/"+estateId+"/stats")
	server := &Server{
		StatsEstateUseCase: &mockStatsEstateUseCaseError{},
	}
	_ = server.GetEstateSummary(ctx, estateId)

	var resp generated.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", resp.Message)

}

func TestServer_GetDistanceForDronePlan_InternalServerError(t *testing.T) {
	estateId := uuid.New().String()
	rec, ctx := createContextAndRecordHttp("", "/estate/"+estateId+"/drone-plan")
	server := &Server{
		StatsEstateUseCase: &mockStatsEstateUseCaseError{},
	}
	_ = server.GetDistanceForDronePlan(ctx, estateId, generated.GetDistanceForDronePlanParams{})

	var resp generated.ErrorResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", resp.Message)

}

type mockStatsEstateUseCaseError struct {
}

func (m *mockStatsEstateUseCaseError) GetDronePlanDistance(ctx context.Context, estateId string, params generated.GetDistanceForDronePlanParams) (generated.GetDronePlanDistance, error) {
	return generated.GetDronePlanDistance{}, errors.New("internal server error")
}

func (m *mockStatsEstateUseCaseError) PublishCalculationStatsEstate(estateId string) {
	//TODO implement me
	panic("implement me")
}

func (m *mockStatsEstateUseCaseError) GetStatsEstate(ctx context.Context, estateId string) (generated.GetStatsEstateResponse, error) {
	return generated.GetStatsEstateResponse{}, errors.New("internal server error")
}

func waitForStatsToBeCalculated(t *testing.T, estateId string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for {
		rec, ctx := createContextAndRecordHttp("", "/estate/"+estateId+"/stats")
		_ = server.GetEstateSummary(ctx, estateId)

		var resp generated.GetStatsEstateResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotNil(t, resp)
		if resp.Max == 5 && resp.Min == 3 && resp.Count == 3 && resp.Median == 4.0 {
			validateDronePlan(t, estateId)
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("Timeout waiting for stats to be calculated")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func validateDronePlan(t *testing.T, estateId string) {
	rec, ctx := createContextAndRecordHttp("", "/estate/"+estateId+"/drone-plan")
	_ = server.GetDistanceForDronePlan(ctx, estateId, generated.GetDistanceForDronePlanParams{})

	var resp generated.GetDronePlanDistance
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.NotNil(t, resp)
	assert.Equal(t, 54, resp.Distance)
}
