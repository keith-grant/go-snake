package grid

import (
	bs "github.com/BattlesnakeOfficial/starter-snake-go/gameTypes"
)

type CellType int32

const (
	EMPTY CellType = 0
	FOOD  CellType = 1
	SNAKE CellType = 2
)

/////////////////////
type Cell struct {
	position           bs.Coord
	Up                 *Cell
	Right              *Cell
	Down               *Cell
	Left               *Cell
	Type               CellType
	DistanceToFood     int
	ConnectedCellCount int
}

func NewCell(x int, y int, cellType CellType) *Cell {
	newCell := new(Cell)
	newCell.position.X = x
	newCell.position.Y = y
	newCell.Up = nil
	newCell.Right = nil
	newCell.Down = nil
	newCell.Left = nil
	newCell.Type = cellType
	return newCell
}

func getCellType(state bs.GameState, x int, y int) CellType {
	retVal := EMPTY
	if isSnake(state, x, y) {
		retVal = SNAKE
	} else if isFood(state, x, y) {
		retVal = FOOD
	}
	return retVal
}

func isSnake(state bs.GameState, x int, y int) bool {
	snakeList := state.Board.Snakes
	for snakeindex := 0; snakeindex < len(snakeList); snakeindex++ {
		snake := snakeList[snakeindex]
		for segmentindex := 0; segmentindex < len(snake.Body); segmentindex++ {
			segment := snake.Body[segmentindex]
			if x == segment.X && y == segment.Y {
				return true
			}
		}
	}
	return false
}

func isFood(state bs.GameState, x int, y int) bool {
	foodList := state.Board.Food
	for foodindex := 0; foodindex < len(foodList); foodindex++ {
		food := foodList[foodindex]
		if x == food.X && y == food.Y {
			return true
		}
	}
	return false
}

func connectToSiblings(cell *Cell, grid *Grid) {
	cell.Up = GetCell(cell.position.X, cell.position.Y+1, grid)
	cell.Down = GetCell(cell.position.X, cell.position.Y-1, grid)
	cell.Left = GetCell(cell.position.X-1, cell.position.Y, grid)
	cell.Right = GetCell(cell.position.X+1, cell.position.Y, grid)
}

func GetCell(x int, y int, grid *Grid) *Cell {

	// check boundaries
	if x < 0 {
		return nil
	}
	if y < 0 {
		return nil
	}
	if x > grid.width-1 {
		return nil
	}
	if y > grid.height-1 {
		return nil
	}

	return grid.allCells[bs.Coord{X: x, Y: y}]
}

//////////////////////////////////

type Grid struct {
	width    int
	height   int
	allCells map[bs.Coord]*Cell
}

func CreateGrid(state bs.GameState) *Grid {
	grid := new(Grid)
	grid.width = state.Board.Width
	grid.height = state.Board.Height
	grid.allCells = make(map[bs.Coord]*Cell)

	// Populate all cells
	for x := 0; x < state.Board.Width; x++ {
		for y := 0; y < state.Board.Height; y++ {
			cellType := getCellType(state, x, y)
			grid.allCells[bs.Coord{X: x, Y: y}] = NewCell(x, y, cellType)
		}
	}

	// connect all cells
	for _, element := range grid.allCells {
		connectToSiblings(element, grid)
	}

	for _, element := range grid.allCells {
		element.DistanceToFood = distanceToFood(*element, state)
	}

	return grid
}

func GetCellPosition(cell *Cell) bs.Coord {
	return cell.position
}

func distanceToFood(cell Cell, state bs.GameState) int {
	foodList := state.Board.Food

	distanceToFood := state.Board.Width + state.Board.Height // max possible distance on a grid
	for entry := range foodList {
		distanceToEntry := distanceBetweenCells(cell.position, foodList[entry])
		if distanceToEntry < distanceToFood {
			distanceToFood = distanceToEntry
		}
	}

	return distanceToFood
}

func distanceBetweenCells(a bs.Coord, b bs.Coord) int {
	return intAbs(a.X-b.X) + intAbs(a.Y-b.Y)
}

func intAbs(x int) int { // No built in integer abs !!!!
	if x < 0 {
		return -x
	}
	return x
}

func isConnected(thisCell *Cell, targetCell *Cell) bool {
	return false
}

type void struct{}

func ConnectedCellCount(cell *Cell) int {

	connectedCellMap := make(map[*Cell]struct{})

	//early out
	_, found := connectedCellMap[cell]
	if found {
		// don't look at this cell again
		return 0
	}
	// For each attached cell add it to the connectedCellMap and then call this on it
	attachedCells := getAttachedCells(cell)
	for entry := range attachedCells {
		var member void
		connectedCellMap[attachedCells[entry]] = member
	}

	// then call connectedCells on all attachedCells
	for entry := range attachedCells {
		connectedCellCountInternal(attachedCells[entry], connectedCellMap)
	}
	return len(connectedCellMap)
}

func connectedCellCountInternal(cell *Cell, connectedCellMap map[*Cell]struct{}) {

	// For each attached cell add it to the connectedCellMap and then call this on it
	attachedCells := getAttachedCells(cell)

	// Keep a note of which attached Cells are being seen for teh first time
	newCells := make(map[*Cell]struct{})
	for entry := range attachedCells {
		_, exists := connectedCellMap[attachedCells[entry]]
		if !exists {
			var member void
			newCells[attachedCells[entry]] = member
		}
	}

	// Add All attached cells to teh set of connected cells
	for entry := range attachedCells {
		var member void
		connectedCellMap[attachedCells[entry]] = member
	}

	// then call connectedCells only on newCells
	for key := range newCells {
		connectedCellCountInternal(key, connectedCellMap)
	}
}

func getAttachedCells(cell *Cell) []*Cell {
	var attachedCells []*Cell

	if validCell(cell.Up) {
		attachedCells = append(attachedCells, cell.Up)
	}
	if validCell(cell.Right) {
		attachedCells = append(attachedCells, cell.Right)
	}
	if validCell(cell.Down) {
		attachedCells = append(attachedCells, cell.Down)
	}
	if validCell(cell.Left) {
		attachedCells = append(attachedCells, cell.Left)
	}

	return attachedCells
}

func validCell(cell *Cell) bool {
	if cell == nil {
		return false
	}
	if cell.Type == SNAKE {
		return false
	}
	return true
}
