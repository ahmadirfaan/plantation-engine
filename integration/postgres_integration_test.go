package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/ahmadirfaan/plantation-engine/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// schema is a trimmed version of database.sql (no CREATE DATABASE / \c prefix,
// no 100-way partitioning) that still exercises the real unique constraints the
// repositories depend on.
const schema = `
CREATE TABLE estate
(
    id         text PRIMARY KEY,
    name       VARCHAR(255),
    width      INT,
    length     INT,
    ext_info   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE block
(
    id           UUID,
    estate_id    UUID NOT NULL,
    x_coordinate INT,
    y_coordinate INT,
    ext_info     TEXT,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, estate_id)
);

CREATE UNIQUE INDEX block_uniq_estate_xy ON block (estate_id, x_coordinate, y_coordinate);

CREATE TABLE tree
(
    id         UUID,
    block_id   UUID NOT NULL,
    estate_id  UUID NOT NULL,
    height     INT,
    ext_info   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, estate_id)
);

CREATE UNIQUE INDEX tree_uniq_block_estate ON tree (block_id, estate_id);

CREATE TABLE estate_stats
(
    estate_id            UUID PRIMARY KEY,
    min_height_tree      DOUBLE PRECISION,
    max_height_tree      DOUBLE PRECISION,
    sum_tree             INT,
    median_height_tree   DOUBLE PRECISION,
    total_distance_drone DOUBLE PRECISION,
    ext_info             TEXT,
    created_at           TIMESTAMP DEFAULT NOW(),
    updated_at           TIMESTAMP DEFAULT NOW(),
    calculated_at        TIMESTAMP DEFAULT NOW()
);
`

func setupTestDB(t *testing.T) *sql.DB {
	if os.Getenv("TESTCONTAINERS") == "" {
		t.Skip("set TESTCONTAINERS=1 to run postgres integration tests")
	}
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "database",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/database?sslmode=disable", port.Port())
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.ExecContext(ctx, schema)
	require.NoError(t, err)

	return db
}

func TestPostgresRepositoriesIntegration(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Run("save and query estate", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		id, err := estateRepo.SaveEstate(ctx, 5, 5)
		require.NoError(t, err)
		require.NotNil(t, id)

		estate, err := estateRepo.QueryByEstateId(ctx, *id)
		require.NoError(t, err)
		require.NotNil(t, estate)
		assert.Equal(t, *id, estate.Id)
		assert.Equal(t, 5, estate.Width)
		assert.Equal(t, 5, estate.Length)
	})

	t.Run("save block and tree atomically", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		estateId, err := estateRepo.SaveEstate(ctx, 3, 3)
		require.NoError(t, err)

		treeRepo := repository.NewTreeRepository(db)
		treeId, err := treeRepo.SaveBlockAndTree(ctx, *estateId, 1, 1, 10)
		require.NoError(t, err)
		require.NotNil(t, treeId)

		blockRepo := repository.NewBlockRepository(db)
		block, err := blockRepo.QueryByEstateIdAndBlockCoordinate(ctx, *estateId, 1, 1)
		require.NoError(t, err)
		require.NotNil(t, block)
		assert.Equal(t, 1, block.BlockX)
		assert.Equal(t, 1, block.BlockY)
	})

	t.Run("duplicate coordinate returns ErrBlockHasTree", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		estateId, err := estateRepo.SaveEstate(ctx, 3, 3)
		require.NoError(t, err)

		treeRepo := repository.NewTreeRepository(db)
		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 2, 2, 10)
		require.NoError(t, err)

		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 2, 2, 20)
		require.ErrorIs(t, err, repository.ErrBlockHasTree)
	})

	t.Run("query all trees joins block coordinates", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		estateId, err := estateRepo.SaveEstate(ctx, 5, 5)
		require.NoError(t, err)

		treeRepo := repository.NewTreeRepository(db)
		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 3, 4, 15)
		require.NoError(t, err)

		statsRepo := repository.NewStatsEstateRepository(db)
		trees, err := statsRepo.QueryAllTree(ctx, *estateId)
		require.NoError(t, err)
		require.Len(t, trees, 1)
		assert.Equal(t, 3, trees[0].XAxis)
		assert.Equal(t, 4, trees[0].YAxis)
		assert.Equal(t, 15, trees[0].Height)
	})

	t.Run("save and query stats roundtrip", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		estateId, err := estateRepo.SaveEstate(ctx, 1, 5)
		require.NoError(t, err)

		treeRepo := repository.NewTreeRepository(db)
		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 2, 1, 5)
		require.NoError(t, err)
		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 3, 1, 3)
		require.NoError(t, err)
		_, err = treeRepo.SaveBlockAndTree(ctx, *estateId, 4, 1, 4)
		require.NoError(t, err)

		statsRepo := repository.NewStatsEstateRepository(db)
		err = statsRepo.SaveStatsEstate(ctx, model.EstateStats{
			EstateId:           *estateId,
			MinHeightTree:      3,
			MaxHeightTree:      5,
			MedianHeightTree:   4.0,
			SumTree:            3,
			TotalDistanceDrone: 54,
		})
		require.NoError(t, err)

		stats, err := statsRepo.QueryById(ctx, *estateId)
		require.NoError(t, err)
		require.NotNil(t, stats)
		assert.Equal(t, 3, stats.MinHeightTree)
		assert.Equal(t, 5, stats.MaxHeightTree)
		assert.Equal(t, 4.0, stats.MedianHeightTree)
		assert.Equal(t, 3, stats.SumTree)
		assert.Equal(t, 54, stats.TotalDistanceDrone)
	})

	t.Run("query missing estate returns nil", func(t *testing.T) {
		estateRepo := repository.NewEstateRepository(db)
		estate, err := estateRepo.QueryByEstateId(ctx, uuid.New().String())
		require.NoError(t, err)
		assert.Nil(t, estate)
	})
}
