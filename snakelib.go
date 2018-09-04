package main

type GameState struct {
	SnakeRequest *SnakeRequest
	Grid         *Grid
}

type Cell struct {
	IsEmpty bool
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

func NewGameState(snakeRequest *SnakeRequest) *GameState {
	grid := emptyGrid(snakeRequest)
	for _, snake := range snakeRequest.Board.Snakes {
		for _, part := range snake.Body {
			grid.Cells[CoordToIndex(part, snakeRequest)].IsEmpty = false
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
