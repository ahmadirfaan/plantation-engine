package helper

import (
	"fmt"
	"testing"

	"github.com/ahmadirfaan/plantation-engine/model"
	"github.com/stretchr/testify/assert"
)

func TestCreatingGrid(t *testing.T) {
	trees := buildBlockData()
	grid := createGrid(model.Estate{
		Width:  3,
		Length: 5,
	}, trees)
	fmt.Println(grid)
	assert.Equal(t, 3, len(grid))
}

func TestTotalDistanceDronePlan(t *testing.T) {
	trees := buildBlockData()
	estateStats := model.EstateStats{}
	CalculateTotalDistanceDrone(trees, model.Estate{
		Width:  1,
		Length: 5,
	}, &estateStats)
	assert.Equal(t, 54, estateStats.TotalDistanceDrone)
}

func TestTotalDistanceDronePlanGrid3x3(t *testing.T) {
	trees := make([]model.Tree, 0)
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  1,
		Height: 10,
	})
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  1,
		Height: 15,
	})
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  2,
		Height: 29,
	})
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  3,
		Height: 30,
	})
	estateStats := model.EstateStats{}
	CalculateTotalDistanceDrone(trees, model.Estate{
		Width:  3,
		Length: 3,
	}, &estateStats)
	assert.Equal(t, 210, estateStats.TotalDistanceDrone)
}

func TestTotalDistanceDronePlanGrid2x2(t *testing.T) {
	trees := make([]model.Tree, 0)
	trees = append(trees, model.Tree{
		XAxis:  1,
		YAxis:  2,
		Height: 10,
	})
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  2,
		Height: 15,
	})
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  1,
		Height: 29,
	})
	estateStats := model.EstateStats{}
	CalculateTotalDistanceDrone(trees, model.Estate{
		Width:  2,
		Length: 2,
	}, &estateStats)
	assert.Equal(t, 90, estateStats.TotalDistanceDrone)
}

func TestLastBlockFinish(t *testing.T) {
	x, y := LastBlockFinish(3, 3)
	assert.Equal(t, 3, x)
	assert.Equal(t, 3, y)

	x, y = LastBlockFinish(3, 2)
	assert.Equal(t, 1, x)
	assert.Equal(t, 2, y)

	x, y = LastBlockFinish(5, 8)
	assert.Equal(t, 1, x)
	assert.Equal(t, 8, y)
}

func TestCheckRestCoordinate(t *testing.T) {
	coordinate := CheckRestCoordinate(40, buildBlockData(), 5, 1)
	assert.Equal(t, 4, coordinate.X)
	assert.Equal(t, 1, coordinate.Y)

	coordinateSecond := CheckRestCoordinate(27, buildBlockData(), 5, 1)
	assert.Equal(t, 2, coordinateSecond.X)
	assert.Equal(t, 1, coordinateSecond.Y)
}

func buildBlockData() []model.Tree {
	trees := make([]model.Tree, 0)
	trees = append(trees, model.Tree{
		XAxis:  2,
		YAxis:  1,
		Height: 5,
	})
	trees = append(trees, model.Tree{
		XAxis:  3,
		YAxis:  1,
		Height: 3,
	})
	trees = append(trees, model.Tree{
		XAxis:  4,
		YAxis:  1,
		Height: 4,
	})
	return trees
}
