package main

import (
	"groupie-tracker/internal/api"
	"log"
	"net/http"
)

const port = ":8088"

func main() {
	mux := http.NewServeMux()

	handler := api.NewHandler()
	handler.SetEndpoints(mux)
	log.Printf("Starting server...\nhttp://localhost%v/\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Println(err)
		return
	}
}
