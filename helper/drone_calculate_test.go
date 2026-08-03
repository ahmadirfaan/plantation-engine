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

func TestCalculateTotalDistanceDrone_TableDriven(t *testing.T) {
	fullGrid := func(width, length, height int) []model.Tree {
		var trees []model.Tree
		for y := 1; y <= width; y++ {
			for x := 1; x <= length; x++ {
				trees = append(trees, model.Tree{XAxis: x, YAxis: y, Height: height})
			}
		}
		return trees
	}

	tests := []struct {
		name   string
		estate model.Estate
		trees  []model.Tree
		want   int
	}{
		{name: "empty estate 3x3", estate: model.Estate{Width: 3, Length: 3}, trees: nil, want: 82},
		{name: "empty estate 1x1", estate: model.Estate{Width: 1, Length: 1}, trees: nil, want: 2},
		{name: "full estate 2x2 height 5", estate: model.Estate{Width: 2, Length: 2}, trees: fullGrid(2, 2, 5), want: 42},
		{name: "full estate 3x3 height 5", estate: model.Estate{Width: 3, Length: 3}, trees: fullGrid(3, 3, 5), want: 92},
		{name: "extreme elevations 30 then 1", estate: model.Estate{Width: 1, Length: 2}, trees: []model.Tree{{XAxis: 1, YAxis: 1, Height: 30}, {XAxis: 2, YAxis: 1, Height: 1}}, want: 72},
		{name: "single tallest tree 1x1", estate: model.Estate{Width: 1, Length: 1}, trees: []model.Tree{{XAxis: 1, YAxis: 1, Height: 30}}, want: 62},
		{name: "known 1x5 case", estate: model.Estate{Width: 1, Length: 5}, trees: buildBlockData(), want: 54},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := model.EstateStats{}
			CalculateTotalDistanceDrone(tt.trees, tt.estate, &stats)
			assert.Equal(t, tt.want, stats.TotalDistanceDrone)
		})
	}
}

func TestCheckRestCoordinate_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		maxDistance int
		wantX       int
		wantY       int
	}{
		{name: "zero distance stays at start", maxDistance: 0, wantX: 1, wantY: 1},
		{name: "short distance", maxDistance: 27, wantX: 2, wantY: 1},
		{name: "mid distance", maxDistance: 40, wantX: 4, wantY: 1},
		{name: "beyond total returns last block", maxDistance: 1000, wantX: 5, wantY: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord := CheckRestCoordinate(tt.maxDistance, buildBlockData(), 5, 1)
			assert.Equal(t, tt.wantX, coord.X)
			assert.Equal(t, tt.wantY, coord.Y)
		})
	}
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
