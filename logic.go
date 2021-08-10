package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
  //"math"
  gb "github.com/BattlesnakeOfficial/starter-snake-go/grid"
  bs "github.com/BattlesnakeOfficial/starter-snake-go/gameTypes"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your Battlesnake username
		Color:      "#888888", // TODO: Personalize
		Head:       "default", // TODO: Personalize
		Tail:       "default", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state bs.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}


// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state bs.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state bs.GameState) BattlesnakeMoveResponse {

  boardGrid := gb.CreateGrid(state)

  up := bs.Coord{X: state.You.Head.X, Y: state.You.Head.Y + 1}
  down := bs.Coord{X: state.You.Head.X, Y: state.You.Head.Y - 1}
  left := bs.Coord{X: state.You.Head.X -1, Y: state.You.Head.Y}
  right := bs.Coord{X: state.You.Head.X + 1, Y: state.You.Head.Y}

	possibleMoves := map[string]*gb.Cell{
		"up":    gb.GetCell(up.X, up.Y, boardGrid),
		"down":  gb.GetCell(down.X, down.Y, boardGrid),
		"left":  gb.GetCell(left.X, left.Y, boardGrid),
		"right": gb.GetCell(right.X, right.Y, boardGrid),
	}
 
  //log.Print(boardGrid)
	// TODO: Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	// boardWidth := state.Board.Width
	// boardHeight := state.Board.Height

//  avoidWalls(state, possibleMoves)

	// TODO: Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake from colliding with itself.
	// mybody := state.You.Body

//  avoidSelf(state, possibleMoves)

	// TODO: Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake from colliding with others.

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := make(map[string]*gb.Cell)
	for move, cell := range possibleMoves {
    if cell != nil{
      //log.Printf("%s : %+v\n", move, cell)
    } else {
      //log.Printf("%s : is nil\n", move)
    }
		if isSafe(cell, state.You.Body) {
			safeMoves[move] = cell
		}
	}


  safeMovesInfo := make(map[string]int)

  for key, entry := range safeMoves{
    entry.ConnectedCellCount = gb.ConnectedCellCount(entry);
    safeMovesInfo[key] = entry.ConnectedCellCount
  }

  log.Print(safeMovesInfo)
	if len(safeMoves) == 0 {
		nextMove = "down"
    //prefer the wall to eating yourself but if there are no walls then default to down
    for move, cell := range possibleMoves {
      if cell == nil {
        nextMove = move
      }
    }
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = decideMove(safeMoves, state)
		log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func isSafe(cell* gb.Cell, me []bs.Coord) bool{
  if cell == nil {
    return false
  }
  switch cell.Type{
    case gb.EMPTY:
    return true
    case gb.FOOD:
    return true
    case gb.SNAKE:
    return false
  }

  return isMe(gb.GetCellPosition(cell).X, gb.GetCellPosition(cell).Y, me)
}

func isMe(targetX int, targetY int, me []bs.Coord) bool {
  for i := 0; i < len(me);  i++ {
    if(me[i].X == targetX && me[i].Y == targetY) {
      return true
    }
  }
  return false
}

func decideMove(options map[string]*gb.Cell, state bs.GameState) string{
  starting_value := 500
  distanceToFood := starting_value
  var move string

  imHungry := hungry(state.You)
  for key, entry := range options{
    // if hungary move towards FOOD
    if(imHungry) {
      if(entry.DistanceToFood < distanceToFood) {
        move = key
        distanceToFood = entry.DistanceToFood
      }
    } else {
      // move away from food
      if(entry.DistanceToFood > distanceToFood || distanceToFood == starting_value) {
        move = key
        distanceToFood = entry.DistanceToFood
      }
    }
  }
  return move
}

func hungry(you bs.Battlesnake) bool {
  if you.Health < 20 {
    return true
  }
  return false
}