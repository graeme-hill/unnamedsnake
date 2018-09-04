package main

import (
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
	dump(decoded)

	state := NewGameState(&decoded)
	direction := FindSafeDirection(decoded.You.Body[0], state)
	respond(res, MoveResponse{
		Move: direction,
	})
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}
