package helper

import "github.com/ahmadirfaan/plantation-engine/model"

func LastBlockFinish(length, width int) (int, int) {
	if width%2 == 0 {
		return 1, width
	}
	return length, width
}

func CalculateTotalDistanceDrone(trees []model.Tree, estate model.Estate, estateStats *model.EstateStats) {
	grid := createGrid(estate, trees)
	total, _ := traverseGrid(grid, 0, false)
	estateStats.TotalDistanceDrone = total
}

func CheckRestCoordinate(maxDistance int, trees []model.Tree, length, width int) model.Coordinate {
	grid := createGrid(model.Estate{Length: length, Width: width}, trees)

	_, rest := traverseGrid(grid, maxDistance, true)
	if rest != nil {
		return *rest
	}

	x, y := LastBlockFinish(length, width)
	return model.Coordinate{X: x, Y: y}
}

func traverseGrid(grid [][]int, maxDistance int, checkRest bool) (int, *model.Coordinate) {
	width, length := len(grid), len(grid[0])
	totalDistance, prevHeight := 0, 0
	var prevX, prevY int
	firstStep := true

	for y := 0; y < width; y++ {
		xRange := getXRange(y, length)
		for _, x := range xRange {
			if !firstStep {
				if checkRest {
					if maxDistance < 10 {
						return totalDistance, &model.Coordinate{X: prevX + 1, Y: prevY + 1}
					}
					maxDistance -= 10
				}
				totalDistance += 10
			}
			treeHeight := grid[y][x]
			droneHeight := getDroneHeight(treeHeight)
			heightDistance := abs(droneHeight - prevHeight)

			if checkRest {
				if maxDistance < heightDistance {
					return totalDistance, &model.Coordinate{X: prevX + 1, Y: prevY + 1}
				}
				maxDistance -= heightDistance
			}

			prevHeight = droneHeight
			totalDistance += heightDistance

			prevX, prevY = x, y
			firstStep = false

		}
	}

	if !checkRest {
		totalDistance += prevHeight
	}

	return totalDistance, nil
}

func getRestCoordinate(x, y int) *model.Coordinate {
	return &model.Coordinate{X: x, Y: y}
}

func createGrid(estate model.Estate, trees []model.Tree) [][]int {
	width, length := estate.Width, estate.Length
	grid := make([][]int, width)
	for i := range grid {
		grid[i] = make([]int, length)
	}

	for _, tree := range trees {
		x := tree.XAxis - 1
		y := tree.YAxis - 1
		if y >= 0 && y < estate.Width && x >= 0 && x < estate.Length {
			grid[y][x] = tree.Height
		}
	}
	return grid
}

func makeRange(start, end, step int) []int {
	var result []int
	for i := start; i != end; i += step {
		result = append(result, i)
	}
	return result
}

func getDroneHeight(treeHeight int) int {
	if treeHeight > 0 {
		return treeHeight + 1
	}
	return 1
}

func getXRange(row, cols int) []int {
	if row%2 == 0 {
		return makeRange(0, cols, 1) // west to east
	}
	return makeRange(cols-1, -1, -1) // east to west
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
