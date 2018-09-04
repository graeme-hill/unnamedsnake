package main

import (
	"fmt"
	"log"
	"net/http"
)

func Start(res http.ResponseWriter, req *http.Request) {
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}
	dump(decoded)

	respond(res, StartResponse{
		Color: "#75CEDD",
	})
}

func Move(res http.ResponseWriter, req *http.Request) {
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	//dump(decoded)

	state := NewGameState(&decoded)
	//direction := FindSafeDirection(decoded.You.Body[0], state)
	destination := ClosestFood(decoded.You.Body[0], state)
	dir := AStar(decoded.You.Body[0], destination, state)
	nonSuicidal := NonSuicidalDirections(decoded.You.Body[0], state)

	if state.SnakeRequest.You.Health > 50 {
		tailDir := DirectionToTail(state)
		fmt.Println("trying tail")
		if tailDir != NOT_FOUND {
			fmt.Println("TAIL")
			dir = tailDir
		}
	}

	dirIsOkay := false
	for _, nonSuidicalDirection := range nonSuicidal {
		if dir == nonSuidicalDirection {
			dirIsOkay = true
			break
		}
	}

	if !dirIsOkay && len(nonSuicidal) > 0 {
		fmt.Println("DON'T KILL YOURSELF!")
		dir = nonSuicidal[0]
	}

	fmt.Println(dir)
	respond(res, MoveResponse{
		Move: dir,
	})
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}
