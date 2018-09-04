package main

const BIG_NUMBER = 9999999
const NOT_FOUND = "nope"

type GameState struct {
	SnakeRequest *SnakeRequest
	Grid         *Grid
}

type Cell struct {
	IsEmpty   bool
	Dangerous bool
}

type Grid struct {
	Cells []Cell
}

func emptyGrid(snakeRequest *SnakeRequest) *Grid {
	var grid Grid
	for x := 0; x < snakeRequest.Board.Width; x++ {
		for y := 0; y < snakeRequest.Board.Height; y++ {
			grid.Cells = append(grid.Cells, Cell{IsEmpty: true})
		}
	}
	return &grid
}

func CoordToIndex(coord Coord, snakeRequest *SnakeRequest) int {
	return coord.Y*snakeRequest.Board.Width + coord.X
}

func IndexToCoord(index int, snakeRequest *SnakeRequest) Coord {
	x := index % snakeRequest.Board.Width
	y := index / snakeRequest.Board.Width
	return Coord{x, y}
}

func Ring(coord Coord) []Coord {
	return []Coord{
		Coord{coord.X, coord.Y - 1},
		Coord{coord.X, coord.Y + 1},
		Coord{coord.X + 1, coord.Y},
		Coord{coord.X - 1, coord.Y},
		Coord{coord.X - 1, coord.Y - 1},
		Coord{coord.X + 1, coord.Y + 1},
		Coord{coord.X + 1, coord.Y - 1},
		Coord{coord.X - 1, coord.Y + 1},
	}
}

func NewGameState(snakeRequest *SnakeRequest) *GameState {
	grid := emptyGrid(snakeRequest)

	// Mark non-empty cells
	for _, snake := range snakeRequest.Board.Snakes {
		for _, part := range snake.Body {
			grid.Cells[CoordToIndex(part, snakeRequest)].IsEmpty = false
		}
	}

	// Mark dangerous cells
	for _, snake := range snakeRequest.Board.Snakes {
		if snake.ID == snakeRequest.You.ID {
			continue
		}
		if len(snake.Body) < len(snakeRequest.You.Body) {
			continue
		}

		for _, cell := range Ring(snake.Body[0]) {
			index := CoordToIndex(cell, snakeRequest)
			grid.Cells[index].Dangerous = true
		}
	}

	return &GameState{SnakeRequest: snakeRequest, Grid: grid}
}

func IsCellSafe(coord Coord, state *GameState) bool {
	if coord.X < 0 || coord.Y < 0 {
		return false
	}

	if coord.X >= state.SnakeRequest.Board.Width || coord.Y >= state.SnakeRequest.Board.Height {
		return false
	}

	return state.Grid.Cells[CoordToIndex(coord, state.SnakeRequest)].IsEmpty
}

func IsCellDangerous(coord Coord, state *GameState) bool {
	if coord.X < 0 || coord.Y < 0 {
		return true
	}

	if coord.X >= state.SnakeRequest.Board.Width || coord.Y >= state.SnakeRequest.Board.Height {
		return true
	}

	return state.Grid.Cells[CoordToIndex(coord, state.SnakeRequest)].Dangerous
}

func FindSafeDirection(start Coord, state *GameState) string {
	if IsCellSafe(Coord{start.X, start.Y - 1}, state) {
		return "up"
	}
	if IsCellSafe(Coord{start.X + 1, start.Y}, state) {
		return "right"
	}
	if IsCellSafe(Coord{start.X, start.Y + 1}, state) {
		return "down"
	}
	return "left"
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ManhattanDistance(a, b Coord) int {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}

func GetNeighbors(coord Coord) []Coord {
	return []Coord{
		Coord{coord.X, coord.Y - 1},
		Coord{coord.X, coord.Y + 1},
		Coord{coord.X + 1, coord.Y},
		Coord{coord.X - 1, coord.Y},
	}
}

func CalcDirection(a, b Coord) string {
	if a.X < b.X {
		return "right"
	}
	if a.X > b.X {
		return "left"
	}
	if a.Y < b.Y {
		return "down"
	}
	return "up"
}

func ClosestFood(start Coord, state *GameState) Coord {
	closest := Coord{0, 0}
	closestDistance := BIG_NUMBER
	for _, food := range state.SnakeRequest.Board.Food {
		distance := ManhattanDistance(start, food)
		if distance < closestDistance {
			closest = food
			closestDistance = distance
		}
	}
	return closest
}

func AStar(start, goal Coord, state *GameState) string {
	closedSet := map[int]interface{}{}
	openSet := map[int]interface{}{}
	cameFrom := map[int]int{}
	gScore := map[int]int{}
	fScore := map[int]int{}

	startIndex := CoordToIndex(start, state.SnakeRequest)
	goalIndex := CoordToIndex(goal, state.SnakeRequest)
	gScore[startIndex] = 0
	fScore[startIndex] = ManhattanDistance(start, goal)
	openSet[startIndex] = nil

	getGScore := func(index int) int {
		score, ok := gScore[index]
		if ok {
			return score
		} else {
			return BIG_NUMBER
		}
	}

	getFScore := func(index int) int {
		score, ok := fScore[index]
		if ok {
			return score
		} else {
			return BIG_NUMBER
		}
	}

	getLowestFScore := func() int {
		lowest := BIG_NUMBER
		best := -1
		for k, _ := range openSet {
			thisFScore := getFScore(k)
			if getFScore(k) < lowest {
				lowest = thisFScore
				best = k
			}
		}
		return best
	}

	reconstructPath := func(index int) string {
		from := index
		lastFrom := index
		for from != startIndex {
			lastFrom = from
			from = cameFrom[from]
		}
		return CalcDirection(IndexToCoord(from, state.SnakeRequest), IndexToCoord(lastFrom, state.SnakeRequest))
	}

	for len(openSet) > 0 {
		currentIndex := getLowestFScore()
		current := IndexToCoord(currentIndex, state.SnakeRequest)
		if currentIndex == goalIndex {
			return reconstructPath(currentIndex)
		}

		delete(openSet, currentIndex)
		closedSet[currentIndex] = nil

		for _, neighbor := range GetNeighbors(current) {
			neighborIndex := CoordToIndex(neighbor, state.SnakeRequest)

			if !IsCellSafe(neighbor, state) {
				continue
			}

			isDangerous := IsCellDangerous(neighbor, state)
			scoreModifier := 0
			if isDangerous {
				scoreModifier += 100
			}

			_, closed := closedSet[neighborIndex]
			if closed {
				continue
			}

			tentativeGScore := getGScore(currentIndex) + ManhattanDistance(current, neighbor) + scoreModifier

			_, open := openSet[neighborIndex]
			if !open {
				openSet[neighborIndex] = nil
			} else if tentativeGScore >= getGScore(neighborIndex) {
				continue
			}

			cameFrom[neighborIndex] = currentIndex
			gScore[neighborIndex] = tentativeGScore
			fScore[neighborIndex] = getGScore(neighborIndex) + ManhattanDistance(neighbor, goal) + scoreModifier
		}
	}

	return NOT_FOUND
}
