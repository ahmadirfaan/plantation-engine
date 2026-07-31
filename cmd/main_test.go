package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	app *echo.Echo
)

func TestNewServerApp(t *testing.T) {
	app = NewServerApp()

	req := httptest.NewRequest(http.MethodGet, "/hello?id=123", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}

func Test_EndToEnd_Happy_Path(t *testing.T) {

	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable")
	body := `{"width":1,"length":5}`
	req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	app = NewServerApp()
	app.ServeHTTP(rec, req)

	// Check if the response ID is returned
	var resp generated.CreatedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, rec.Code, http.StatusCreated)
	assert.NoError(t, err)
	estateId := resp.Id
	assert.NotEmpty(t, estateId)

	var wg sync.WaitGroup
	trees := buildTreesData()
	for _, tree := range trees {
		wg.Add(1)
		go func(x, y, h int) {
			defer wg.Done()
			addToEstate(t, x, y, h, estateId)
		}(tree[0], tree[1], tree[2])
	}

	wg.Wait()
	sumTree := 3
	minHeight := 3
	maxHeight := 5
	totalDistance := 54
	medianHeight := 4.0
	expectedStats := ExpectedStats{
		SumTree:       &sumTree,
		MinHeight:     &minHeight,
		MaxHeight:     &maxHeight,
		MedianHeight:  &medianHeight,
		TotalDistance: &totalDistance,
	}
	waitForStatsToBeCalculated(t, estateId, 10*time.Second, &expectedStats)

	getEstateSummary(t, estateId, &expectedStats)

}

func TestGetStatsEstateShowErrorNotFound(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable")

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/stats", uuid.New().String()), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	app = NewServerApp()
	app.ServeHTTP(rec, req)

	var resp generated.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, rec.Code, http.StatusNotFound)
	assert.Equal(t, "estate not found", resp.Message)
}

func TestGetDronePlanShowErrorNotFound(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable")

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/drone-plan", uuid.New().String()), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	app = NewServerApp()
	app.ServeHTTP(rec, req)

	var resp generated.ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, rec.Code, http.StatusNotFound)
	assert.Equal(t, "estate not found", resp.Message)
}

func getEstateSummary(t *testing.T, estateId string, stats *ExpectedStats) {
	req := httptest.NewRequest(http.MethodGet, "/estate/"+estateId+"/stats", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	// Check if the response ID is returned
	var respSummary generated.GetStatsEstateResponse
	err := json.Unmarshal(rec.Body.Bytes(), &respSummary)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.NoError(t, err)
	assert.Equal(t, *stats.MaxHeight, respSummary.Max)
	assert.Equal(t, *stats.MinHeight, respSummary.Min)
	assert.Equal(t, *stats.SumTree, respSummary.Count)
	assert.Equal(t, *stats.MedianHeight, float64(respSummary.Median))
}

func buildTreesData() [][]int {
	trees := make([][]int, 0)
	firstTree := make([]int, 0)
	firstTree = append(firstTree, 2)
	firstTree = append(firstTree, 1)
	firstTree = append(firstTree, 5)
	trees = append(trees, firstTree)
	secondTree := make([]int, 0)
	secondTree = append(secondTree, 3)
	secondTree = append(secondTree, 1)
	secondTree = append(secondTree, 3)
	trees = append(trees, secondTree)
	thirdTree := make([]int, 0)
	thirdTree = append(thirdTree, 4)
	thirdTree = append(thirdTree, 1)
	thirdTree = append(thirdTree, 4)
	trees = append(trees, thirdTree)
	return trees
}

func addToEstate(t *testing.T, x int, y int, height int, estateId string) {
	body := fmt.Sprintf(`{"x":%d,"y":%d,"height":%d}`, x, y, height)
	req := httptest.NewRequest(http.MethodPost, "/estate/"+estateId+"/tree", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	// Check if the response ID is returned
	var respTree generated.CreatedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &respTree)
	assert.Equal(t, rec.Code, http.StatusCreated)
	assert.NoError(t, err)
	assert.NotEmpty(t, respTree.Id)
}

func waitForStatsToBeCalculated(t *testing.T, estateId string, timeout time.Duration, expectedStats *ExpectedStats) {
	dbDsn := os.Getenv("DATABASE_URL")
	db, _ := sql.Open("postgres", dbDsn)
	deadline := time.Now().Add(timeout)
	for {
		actualStats := ExpectedStats{}
		err := db.QueryRow("SELECT sum_tree, min_height_tree, max_height_tree, total_distance_drone, median_height_tree FROM estate_stats WHERE estate_id = $1", estateId).
			Scan(&actualStats.SumTree, &actualStats.MinHeight, &actualStats.MaxHeight, &actualStats.TotalDistance, &actualStats.MedianHeight)
		if err != nil && err != sql.ErrNoRows {
			t.Fatal("Failed to query estate_stats:", err)
		}
		if statsEqual(expectedStats, &actualStats) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("Timeout waiting for stats to be calculated")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

type ExpectedStats struct {
	SumTree       *int
	MinHeight     *int
	MaxHeight     *int
	TotalDistance *int
	MedianHeight  *float64
}

func statsEqual(a, b *ExpectedStats) bool {
	return a.SumTree != nil && b.SumTree != nil && *a.SumTree == *b.SumTree &&
		a.MinHeight != nil && b.MinHeight != nil && *a.MinHeight == *b.MinHeight &&
		a.MaxHeight != nil && b.MaxHeight != nil && *a.MaxHeight == *b.MaxHeight &&
		a.TotalDistance != nil && b.TotalDistance != nil && *a.TotalDistance == *b.TotalDistance &&
		a.MedianHeight != nil && b.MedianHeight != nil && *a.MedianHeight == *b.MedianHeight
}
